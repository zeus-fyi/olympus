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
		aggCycle := wfExecParams.CycleCountTaskRelative.AggNormalizedCycleCounts[*aggInst.AggTaskID]
		if i%aggCycle == 0 {
			if aggInst.AggTaskID == nil || md.AggregateAnalysis[*aggInst.AggTaskID] == nil {
				continue
			}
			if md.AggregateAnalysis[*aggInst.AggTaskID][aggInst.AnalysisTaskID] == false {
				continue
			}
			retrievalCtx := workflow.WithActivityOptions(ctx, ao)
			window := artemis_orchestrations.CalculateTimeWindowFromCycles(wfExecParams.WorkflowExecTimekeepingParams.RunWindow.UnixStartTime, i-aggCycle, i, wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize)
			var dataIn []artemis_orchestrations.AIWorkflowAnalysisResult
			depM := artemis_orchestrations.MapDependencies(wfExecParams.WorkflowTasks)
			var analysisDep []int
			for k, _ := range depM.AggregateAnalysis[*aggInst.AggTaskID] {
				analysisDep = append(analysisDep, k)
			}
			err := workflow.ExecuteActivity(retrievalCtx, z.AiAggregateAnalysisRetrievalTask, window, []int{oj.OrchestrationID}, analysisDep).Get(retrievalCtx, &dataIn)
			if err != nil {
				logger.Error("failed to run aggregate retrieval", "Error", err)
				return err
			}
			md.AggregateAnalysis[*aggInst.AggTaskID][aggInst.AnalysisTaskID] = false
			if len(dataIn) == 0 {
				continue
			}
			aggCtx := workflow.WithActivityOptions(ctx, ao)
			var aiAggResp *ChatCompletionQueryResponse
			if aggInst.AggTokenOverflowStrategy == nil || aggInst.AggModel == nil || aggInst.AggTaskID == nil {
				continue
			}
			pr := &PromptReduction{
				TokenOverflowStrategy: aws.StringValue(aggInst.AggTokenOverflowStrategy),
				Model:                 aws.StringValue(aggInst.AggModel),
			}
			var sg *hera_search.SearchResultGroup
			if cp.AnalysisEvalActionParams != nil && cp.AnalysisEvalActionParams.SearchResultGroup != nil {
				sg = cp.AnalysisEvalActionParams.SearchResultGroup
			}
			if sg != nil && len(sg.SearchResults) > 0 {
				pr.PromptReductionSearchResults = &PromptReductionSearchResults{
					InSearchGroup: sg,
				}
				if aggInst.AggPrompt != nil {
					pr.PromptReductionText = &PromptReductionText{
						InPromptBody: aws.StringValue(aggInst.AggPrompt),
					}
				}
			} else {
				if aggInst.AggPrompt == nil {
					continue
				}
				pr.PromptReductionText = &PromptReductionText{
					InPromptBody: aws.StringValue(aggInst.AggPrompt),
				}
			}
			wr := &artemis_orchestrations.AIWorkflowAnalysisResult{
				OrchestrationID:       oj.OrchestrationID,
				SourceTaskID:          aws.IntValue(aggInst.AggTaskID),
				RunningCycleNumber:    i,
				SearchWindowUnixStart: window.UnixStartTime,
				SearchWindowUnixEnd:   window.UnixEndTime,
			}

			chunkedTaskCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(chunkedTaskCtx, z.TokenOverflowReduction, ou, pr).Get(chunkedTaskCtx, &pr)
			if err != nil {
				logger.Error("failed to run analysis json", "Error", err)
				return err
			}
			chunkIterator := 0
			if pr.PromptReductionSearchResults != nil && pr.PromptReductionSearchResults.OutSearchGroups != nil {
				chunkIterator = len(pr.PromptReductionSearchResults.OutSearchGroups)
			}
			if pr.PromptReductionText.OutPromptChunks != nil && len(pr.PromptReductionText.OutPromptChunks) > chunkIterator {
				chunkIterator = len(pr.PromptReductionText.OutPromptChunks)
			}
			for chunkOffset := 0; chunkOffset < chunkIterator; chunkOffset++ {
				if pr.PromptReductionSearchResults != nil && pr.PromptReductionSearchResults.OutSearchGroups != nil && chunkOffset < len(pr.PromptReductionSearchResults.OutSearchGroups) {
					sg = pr.PromptReductionSearchResults.OutSearchGroups[chunkOffset]
				} else {
					sg = &hera_search.SearchResultGroup{
						SearchResults: []hera_search.SearchResult{},
					}
				}
				if pr.PromptReductionText.OutPromptChunks != nil && chunkOffset < len(pr.PromptReductionText.OutPromptChunks) {
					aggInst.AggPrompt = &pr.PromptReductionText.OutPromptChunks[chunkOffset]
				}
				wr.ChunkOffset = chunkOffset
				sg.ExtractionPromptExt = aws.StringValue(aggInst.AggPrompt)
				sg.SourceTaskID = aws.IntValue(aggInst.AggTaskID)
				tte := TaskToExecute{
					Ou:  ou,
					Wft: aggInst,
					Sg:  sg,
					Wr:  wr,
				}
				wfID := oj.OrchestrationName + "-agg-json-task-" + strconv.Itoa(i)
				tte.WfID = wfID
				if aggInst.AggTaskName == nil || aggInst.AggModel == nil || aggInst.AggTokenOverflowStrategy == nil || aggInst.AggPrompt == nil || aggInst.AggTaskID == nil {
					return nil
				}
				pr.Model = aws.StringValue(tte.Wft.AggModel)
				pr.TokenOverflowStrategy = aws.StringValue(tte.Wft.AggTokenOverflowStrategy)
				pr.PromptReductionText.InPromptBody = aws.StringValue(tte.Wft.AggPrompt)
				switch aws.StringValue(aggInst.AggResponseFormat) {
				case jsonFormat:
					childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
						WorkflowID:               oj.OrchestrationName + "-agg-json-task-" + strconv.Itoa(i),
						WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
					}
					tte.Tc = TaskContext{
						TaskName:       aws.StringValue(aggInst.AggTaskName),
						TaskType:       AggTask,
						ResponseFormat: aws.StringValue(aggInst.AggResponseFormat),
						Model:          aws.StringValue(aggInst.AggModel),
						TaskID:         aws.IntValue(aggInst.AggTaskID),
					}
					var fullTaskDef []artemis_orchestrations.AITaskLibrary
					selectTaskCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(selectTaskCtx, z.SelectTaskDefinition, tte.Ou, tte.Sg.SourceTaskID).Get(selectTaskCtx, &fullTaskDef)
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
					sg.SearchResults = aiAggResp.FilteredSearchResults
				case socialMediaEngagementResponseFormat:
					if cp.AnalysisEvalActionParams != nil && cp.AnalysisEvalActionParams.SearchResultGroup != nil {
						sg = cp.AnalysisEvalActionParams.SearchResultGroup
					}
					if sg == nil || len(sg.SearchResults) == 0 {
						continue
					}
					sg.ExtractionPromptExt = aws.StringValue(aggInst.AggPrompt)
					childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
						WorkflowID:               oj.OrchestrationName + "-agg-social-media-engagement-" + strconv.Itoa(i),
						WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
					}
					childAggWfCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
					err = workflow.ExecuteChildWorkflow(childAggWfCtx, z.SocialMediaEngagementWorkflow, ou, sg).Get(childAggWfCtx, &aiAggResp)
					if err != nil {
						logger.Error("failed to execute child social media engagement workflow", "Error", err)
						return err
					}
				default:
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
					recordAggCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(recordAggCtx, z.SaveTaskOutput, wr, dataIn).Get(recordAggCtx, nil)
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
