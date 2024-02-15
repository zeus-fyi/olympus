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
	runCycle := cp.RunCycle

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
		if aggInst.AggTaskName == nil || aggInst.AggModel == nil || aggInst.AggTokenOverflowStrategy == nil {
			return nil
		}
		if md.AggregateAnalysis[*aggInst.AggTaskID] == nil {
			continue
		}
		if md.AggregateAnalysis[*aggInst.AggTaskID][aggInst.AnalysisTaskID] == false {
			continue
		}
		aggCycle := wfExecParams.CycleCountTaskRelative.AggNormalizedCycleCounts[*aggInst.AggTaskID]
		if i%aggCycle == 0 {
			window := artemis_orchestrations.CalculateTimeWindowFromCycles(wfExecParams.WorkflowExecTimekeepingParams.RunWindow.UnixStartTime, i-aggCycle, i, wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize)
			depM := artemis_orchestrations.MapDependencies(wfExecParams.WorkflowTasks)
			var analysisDep []int
			for k, _ := range depM.AggregateAnalysis[*aggInst.AggTaskID] {
				analysisDep = append(analysisDep, k)
			}

			var aggRet AggRetResp
			aggRetrievalCtx := workflow.WithActivityOptions(ctx, ao)
			err := workflow.ExecuteActivity(aggRetrievalCtx, z.AiAggregateAnalysisRetrievalTask, window, []int{oj.OrchestrationID}, analysisDep).Get(aggRetrievalCtx, &aggRet)
			if err != nil {
				logger.Error("failed to run aggregate retrieval", "Error", err)
				return err
			}
			md.AggregateAnalysis[*aggInst.AggTaskID][aggInst.AnalysisTaskID] = false
			wr := &artemis_orchestrations.AIWorkflowAnalysisResult{
				OrchestrationID:       oj.OrchestrationID,
				SourceTaskID:          aws.IntValue(aggInst.AggTaskID),
				RunningCycleNumber:    i,
				SearchWindowUnixStart: window.UnixStartTime,
				SearchWindowUnixEnd:   window.UnixEndTime,
			}
			tte := TaskToExecute{
				Ou:  ou,
				Wft: aggInst,
				Sg: &hera_search.SearchResultGroup{
					SourceTaskID:   aws.IntValue(aggInst.AggTaskID),
					Model:          aws.StringValue(aggInst.AggModel),
					ResponseFormat: aws.StringValue(aggInst.AggResponseFormat),
				},
				Wr: wr,
			}

			pr := &PromptReduction{
				Model:                     aws.StringValue(aggInst.AggModel),
				TokenOverflowStrategy:     aws.StringValue(aggInst.AggTokenOverflowStrategy),
				MarginBuffer:              aws.Float64Value(aggInst.AggMarginBuffer),
				DataInAnalysisAggregation: aggRet.InputDataAnalysisToAggSlice,
				AIWorkflowAnalysisResults: aggRet.AIWorkflowAnalysisResultSlice,
				PromptReductionSearchResults: &PromptReductionSearchResults{
					InPromptBody: aws.StringValue(aggInst.AggPrompt),
				},
				PromptReductionText: &PromptReductionText{
					InPromptSystem: aws.StringValue(aggInst.AggPrompt),
				},
			}

			// todo, maybe can deprecate tte.Sg.FilteredSearchResults == nil
			if len(aggRet.InputDataAnalysisToAggSlice) == 0 && len(aggRet.AIWorkflowAnalysisResultSlice) == 0 && tte.Sg.FilteredSearchResults == nil {
				logger.Info("no data in for agg", "aggInst", aggInst)
				continue
			}
			chunkedTaskCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(chunkedTaskCtx, z.TokenOverflowReduction, ou, pr).Get(chunkedTaskCtx, &pr)
			if err != nil {
				logger.Error("failed to run agg token overflow reduction task", "Error", err)
				return err
			}
			chunkIterator := getChunkIteratorLen(pr)
			tte.Tc = TaskContext{
				TaskName:       aws.StringValue(aggInst.AggTaskName),
				TaskType:       AggTask,
				ResponseFormat: aws.StringValue(aggInst.AggResponseFormat),
				Model:          aws.StringValue(aggInst.AggModel),
				TaskID:         aws.IntValue(aggInst.AggTaskID),
			}
			for chunkOffset := 0; chunkOffset < chunkIterator; chunkOffset++ {
				wr.ChunkOffset = chunkOffset
				tte.WfID = oj.OrchestrationName + "-agg-json-task-" + strconv.Itoa(i) + "-" + strconv.Itoa(chunkOffset)
				var sg *hera_search.SearchResultGroup
				if pr.PromptReductionSearchResults != nil && pr.PromptReductionSearchResults.OutSearchGroups != nil && chunkOffset < len(pr.PromptReductionSearchResults.OutSearchGroups) {
					sg = pr.PromptReductionSearchResults.OutSearchGroups[chunkOffset]
					sg.Model = aws.StringValue(aggInst.AggModel)
					sg.ResponseFormat = aws.StringValue(aggInst.AggResponseFormat)
					if chunkOffset < len(pr.PromptReductionText.OutPromptChunks) {
						sg.BodyPrompt = pr.PromptReductionText.OutPromptChunks[chunkOffset]
					}
				} else {
					sg = &hera_search.SearchResultGroup{
						Model:          aws.StringValue(aggInst.AggModel),
						ResponseFormat: aws.StringValue(aggInst.AggResponseFormat),
						BodyPrompt:     pr.PromptReductionText.OutPromptChunks[chunkOffset],
						SearchResults:  []hera_search.SearchResult{},
					}
				}
				tte.Sg = sg
				var aiAggResp *ChatCompletionQueryResponse
				switch aws.StringValue(aggInst.AggResponseFormat) {
				case jsonFormat:
					childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
						WorkflowID:               oj.OrchestrationName + "-agg-json-task-" + strconv.Itoa(i),
						WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
						RetryPolicy:              ao.RetryPolicy,
					}
					var fullTaskDef []artemis_orchestrations.AITaskLibrary
					selectTaskCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(selectTaskCtx, z.SelectTaskDefinition, tte.Ou, aws.IntValue(aggInst.AggTaskID)).Get(selectTaskCtx, &fullTaskDef)
					if err != nil {
						logger.Error("failed to run task", "Error", err)
						return err
					}
					if len(fullTaskDef) == 0 {
						return nil
					}
					var jdef []*artemis_orchestrations.JsonSchemaDefinition
					for _, taskDef := range fullTaskDef {
						jdef = append(jdef, taskDef.Schemas...)
					}
					tte.Tc.Schemas = jdef
					if aggInst.AggTemperature != nil {
						tte.Tc.Temperature = float32(*aggInst.AggTemperature)
					} else {
						tte.Tc.Temperature = 0.5
					}
					childAggWfCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
					err = workflow.ExecuteChildWorkflow(childAggWfCtx, z.JsonOutputTaskWorkflow, tte).Get(childAggWfCtx, &aiAggResp)
					if err != nil {
						logger.Error("failed to execute json agg workflow", "Error", err)
						return err
					}
					wr.WorkflowResultID = aiAggResp.WorkflowResultID
				default:
					aggCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(aggCtx, z.AiAggregateTask, ou, aggInst, tte.Sg).Get(aggCtx, &aiAggResp)
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
					wr.ResponseID = aggRespId
					recordAggCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(recordAggCtx, z.SaveTaskOutput, wr, aiAggResp).Get(recordAggCtx, &aiAggResp.WorkflowResultID)
					if err != nil {
						logger.Error("failed to save aggregation resp", "Error", err)
						return err
					}
				}
				if aiAggResp == nil || len(aiAggResp.Response.Choices) == 0 {
					continue
				}
				if aggInst.AggEvalFns == nil || len(aggInst.AggEvalFns) == 0 {
					continue
				}
				evalWfID := oj.OrchestrationName + "-agg-eval-" + strconv.Itoa(i)
				childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
					WorkflowID:               oj.OrchestrationName + "-agg-eval-" + strconv.Itoa(i),
					RetryPolicy:              ao.RetryPolicy,
					WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
				}
				cp.Window = window
				cp.WfID = evalWfID
				cp.WorkflowResult = *wr
				ea := &EvalActionParams{
					WorkflowTemplateData: aggInst,
					ParentOutputToEval:   aiAggResp,
					EvalFns:              aggInst.AggEvalFns,
					SearchResultGroup:    tte.Sg,
					TaskToExecute:        tte,
				}
				for _, evalFn := range ea.EvalFns {
					evalAggCycle := wfExecParams.CycleCountTaskRelative.AggEvalNormalizedCycleCounts[*aggInst.AggTaskID][evalFn.EvalID]
					if i%evalAggCycle == 0 {
						childAggEvalWfCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
						err = workflow.ExecuteChildWorkflow(childAggEvalWfCtx, z.RunAiWorkflowAutoEvalProcess, cp, ea).Get(childAggEvalWfCtx, nil)
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
