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
		StartToCloseTimeout: time.Minute * 15, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    10,
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
			var dataIn []artemis_orchestrations.AIWorkflowAnalysisResult
			depM := artemis_orchestrations.MapDependencies(wfExecParams.WorkflowTasks)
			var analysisDep []int
			for k, _ := range depM.AggregateAnalysis[*aggInst.AggTaskID] {
				analysisDep = append(analysisDep, k)
			}
			aggRetrievalCtx := workflow.WithActivityOptions(ctx, ao)
			err := workflow.ExecuteActivity(aggRetrievalCtx, z.AiAggregateAnalysisRetrievalTask, window, []int{oj.OrchestrationID}, analysisDep).Get(aggRetrievalCtx, &dataIn)
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
					DataIn:                dataIn,
					SourceTaskID:          aws.IntValue(aggInst.AggTaskID),
					Model:                 aws.StringValue(aggInst.AggModel),
					ResponseFormat:        aws.StringValue(aggInst.AggResponseFormat),
					FilteredSearchResults: cp.GetFilteredSearchResults(),
				},
				Wr: wr,
			}
			pr := &PromptReduction{
				Model:                     aws.StringValue(aggInst.AggModel),
				TokenOverflowStrategy:     aws.StringValue(aggInst.AggTokenOverflowStrategy),
				DataInAnalysisAggregation: dataIn,
				PromptReductionSearchResults: &PromptReductionSearchResults{
					InPromptBody:  aws.StringValue(aggInst.AggPrompt),
					InSearchGroup: tte.Sg,
				},
				PromptReductionText: &PromptReductionText{
					InPromptSystem: aws.StringValue(aggInst.AggPrompt),
				},
			}
			if len(dataIn) == 0 && tte.Sg.FilteredSearchResults == nil {
				logger.Info("no data in for agg", "aggInst", aggInst)
				continue
			}
			var aiAggResp *ChatCompletionQueryResponse
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
				if pr.PromptReductionText.OutPromptChunks != nil && chunkOffset < len(pr.PromptReductionText.OutPromptChunks) {
					tte.Sg.BodyPrompt = pr.PromptReductionText.OutPromptChunks[chunkOffset]
				}
				switch aws.StringValue(aggInst.AggResponseFormat) {
				case jsonFormat:
					childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
						WorkflowID:               oj.OrchestrationName + "-agg-json-task-" + strconv.Itoa(i),
						WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
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
					childAggWfCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
					err = workflow.ExecuteChildWorkflow(childAggWfCtx, z.JsonOutputTaskWorkflow, tte).Get(childAggWfCtx, &aiAggResp)
					if err != nil {
						logger.Error("failed to execute json agg workflow", "Error", err)
						return err
					}
					tte.Sg.FilteredSearchResults = aiAggResp.FilteredSearchResults
					wr.WorkflowResultID = aiAggResp.WorkflowResultID
				case socialMediaEngagementResponseFormat:
					//if cp.AnalysisEvalActionParams != nil && cp.AnalysisEvalActionParams.SearchResultGroup != nil {
					//	sg = cp.AnalysisEvalActionParams.SearchResultGroup
					//}
					//if sg == nil || len(sg.SearchResults) == 0 {
					//	continue
					//}
					//sg.ExtractionPromptExt = aws.StringValue(aggInst.AggPrompt)
					//childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
					//	WorkflowID:               oj.OrchestrationName + "-agg-social-media-engagement-" + strconv.Itoa(i),
					//	WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
					//}
					//childAggWfCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
					//err = workflow.ExecuteChildWorkflow(childAggWfCtx, z.SocialMediaEngagementWorkflow, ou, sg).Get(childAggWfCtx, &aiAggResp)
					//if err != nil {
					//	logger.Error("failed to execute child social media engagement workflow", "Error", err)
					//	return err
					//}
				default:
					aggCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(aggCtx, z.AiAggregateTask, ou, aggInst, dataIn).Get(aggCtx, &aiAggResp)
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
					err = workflow.ExecuteActivity(recordAggCtx, z.SaveTaskOutput, wr, dataIn).Get(recordAggCtx, &aiAggResp.WorkflowResultID)
					if err != nil {
						logger.Error("failed to save aggregation resp", "Error", err)
						return err
					}
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
