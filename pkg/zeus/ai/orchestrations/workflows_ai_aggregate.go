package ai_platform_service_orchestrations

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type InputDataAnalysisToAgg struct {
	ChatCompletionQueryResponse *ChatCompletionQueryResponse   `json:"chatCompletionQueryResponse,omitempty"`
	SearchResultGroup           *hera_search.SearchResultGroup `json:"baseSearchResultsGroup,omitempty"`
}

func (z *ZeusAiPlatformServiceWorkflows) RunAiChildAggAnalysisProcessWorkflow(ctx workflow.Context, cp *MbChildSubProcessParams) error {
	if cp == nil {
		return nil
	}
	wfExecParams := cp.WfExecParams
	ou := cp.Ou
	oj := cp.Oj
	runCycle := cp.Wsr.RunCycle

	md := artemis_orchestrations.MapDependencies(wfExecParams.WorkflowTasks)
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    25,
		},
	}
	i := runCycle
	for _, aggInst := range wfExecParams.WorkflowTasks {
		if aggInst.AggTaskID == nil || aggInst.AggCycleCount == nil || aggInst.AggPrompt == nil || aggInst.AggModel == nil || wfExecParams.WorkflowTaskRelationships.AggAnalysisTasks == nil {
			continue
		}
		if aggInst.AggTaskName == nil || aggInst.AggModel == nil || aggInst.AggTokenOverflowStrategy == nil || md.AggregateAnalysis[*aggInst.AggTaskID] == nil || md.AggregateAnalysis[*aggInst.AggTaskID][aggInst.AnalysisTaskID] == false {
			return nil
		}
		aggCycle := wfExecParams.CycleCountTaskRelative.AggNormalizedCycleCounts[*aggInst.AggTaskID]
		if i%aggCycle == 0 {
			window := artemis_orchestrations.CalculateTimeWindowFromCycles(wfExecParams.WorkflowExecTimekeepingParams.RunWindow.UnixStartTime, i-aggCycle, i, wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize)
			cp.Window = window
			cp.Tc = TaskContext{
				TaskType:              AggTask,
				Model:                 aws.StringValue(aggInst.AggModel),
				TaskID:                aws.IntValue(aggInst.AggTaskID),
				ResponseFormat:        aws.StringValue(aggInst.AggResponseFormat),
				Prompt:                aws.StringValue(aggInst.AggPrompt),
				WorkflowTemplateData:  aggInst,
				TokenOverflowStrategy: aws.StringValue(aggInst.AggTokenOverflowStrategy),
				MarginBuffer:          aws.Float64Value(aggInst.AggMarginBuffer),
				Temperature:           float32(aws.Float64Value(aggInst.AggTemperature)),
			}
			depM := artemis_orchestrations.MapDependencies(wfExecParams.WorkflowTasks)
			var analysisDep []int
			for k, _ := range depM.AggregateAnalysis[*aggInst.AggTaskID] {
				analysisDep = append(analysisDep, k)
			}
			aggRetrievalCtx := workflow.WithActivityOptions(ctx, ao)
			err := workflow.ExecuteActivity(aggRetrievalCtx, z.AiAggregateAnalysisRetrievalTask, cp, analysisDep).Get(aggRetrievalCtx, nil)
			if err != nil {
				logger.Error("failed to run aggregate retrieval", "Error", err)
				return err
			}
			md.AggregateAnalysis[*aggInst.AggTaskID][aggInst.AnalysisTaskID] = false
			var chunkIterator int
			chunkedTaskCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(chunkedTaskCtx, z.TokenOverflowReduction, cp.Wsr.InputID).Get(chunkedTaskCtx, &chunkIterator)
			if err != nil {
				logger.Error("failed to run agg token overflow reduction task", "Error", err)
				return err
			}
			for chunkOffset := 0; chunkOffset < chunkIterator; chunkOffset++ {
				logger.Info("RunAiChildAggAnalysisProcessWorkflow: chunkOffset", chunkOffset)
				cp.Wsr.ChunkOffset = chunkOffset
				switch aws.StringValue(aggInst.AggResponseFormat) {
				case jsonFormat:
					childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{WorkflowID: oj.OrchestrationName + "-agg-json-task-" + strconv.Itoa(i), WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize, RetryPolicy: ao.RetryPolicy}
					childAggWfCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
					cp.Wsr.ChildWfID = childAnalysisWorkflowOptions.WorkflowID
					err = workflow.ExecuteChildWorkflow(childAggWfCtx, z.JsonOutputTaskWorkflow, cp).Get(childAggWfCtx, &cp)
					if err != nil {
						logger.Error("failed to execute json agg workflow", "Error", err)
						return err
					}
				default:
					var aiAggResp *ChatCompletionQueryResponse
					aggCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(aggCtx, z.AiAggregateTask, ou, aggInst, cp).Get(aggCtx, &aiAggResp)
					if err != nil {
						logger.Error("failed to run aggregation", "Error", err)
						return err
					}
					var aggRespId int
					aggCompCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(aggCompCtx, z.RecordCompletionResponse, ou, aiAggResp).Get(aggCompCtx, &aggRespId)
					if err != nil {
						logger.Error("failed to save agg response", "Error", err)
						return err
					}
					wr := &artemis_orchestrations.AIWorkflowAnalysisResult{
						OrchestrationID:       cp.Oj.OrchestrationID,
						SourceTaskID:          cp.Tc.TaskID,
						IterationCount:        0,
						ChunkOffset:           chunkOffset,
						RunningCycleNumber:    cp.Wsr.RunCycle,
						SearchWindowUnixStart: cp.Window.UnixStartTime,
						SearchWindowUnixEnd:   cp.Window.UnixEndTime,
						ResponseID:            aggRespId,
					}
					ia := InputDataAnalysisToAgg{
						ChatCompletionQueryResponse: aiAggResp,
					}
					recordAggCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(recordAggCtx, z.SaveTaskOutput, wr, cp, ia).Get(recordAggCtx, &aiAggResp.WorkflowResultID)
					if err != nil {
						logger.Error("failed to save aggregation resp", "Error", err)
						return err
					}
				}
				for ind, evalFn := range aggInst.AggEvalFns {
					evalAggCycle := wfExecParams.CycleCountTaskRelative.AggEvalNormalizedCycleCounts[*aggInst.AggTaskID][evalFn.EvalID]
					if i%evalAggCycle == 0 {
						childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
							WorkflowID:               oj.OrchestrationName + "-agg-eval-" + strconv.Itoa(i) + "-chunk-" + strconv.Itoa(chunkOffset) + "-ind-" + strconv.Itoa(ind),
							RetryPolicy:              ao.RetryPolicy,
							WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
						}
						cp.Tc.EvalID = evalFn.EvalID
						childAggEvalWfCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
						err = workflow.ExecuteChildWorkflow(childAggEvalWfCtx, z.RunAiWorkflowAutoEvalProcess, cp).Get(childAggEvalWfCtx, nil)
						if err != nil {
							logger.Error("failed to execute child agg eval workflow", "Error", err)
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func getChunkIteratorLen(pr *PromptReduction) int {
	if pr == nil {
		return 0
	}
	chunkIterator := 0
	if pr.PromptReductionSearchResults != nil && pr.PromptReductionSearchResults.OutSearchGroups != nil {
		return len(pr.PromptReductionSearchResults.OutSearchGroups)
	}
	if pr.PromptReductionText != nil && pr.PromptReductionText.OutPromptChunks != nil && len(pr.PromptReductionText.OutPromptChunks) > chunkIterator {
		return len(pr.PromptReductionText.OutPromptChunks)
	}
	return chunkIterator
}
