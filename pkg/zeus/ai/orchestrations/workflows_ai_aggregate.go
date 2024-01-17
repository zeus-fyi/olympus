package ai_platform_service_orchestrations

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (z *ZeusAiPlatformServiceWorkflows) RunAiWorkflowChildAggAnalysisProcess(ctx workflow.Context, cp *MbChildSubProcessParams) error {
	wfExecParams := cp.WfExecParams
	ou := cp.Ou
	oj := cp.Oj
	runCycle := cp.RunCycle

	md := artemis_orchestrations.MapDependencies(wfExecParams.WorkflowTasks)
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15, // Setting a valid non-zero timeout
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
			aoAiAct := workflow.ActivityOptions{
				StartToCloseTimeout: time.Minute * 15, // Setting a valid non-zero timeout
				RetryPolicy: &temporal.RetryPolicy{
					InitialInterval:    time.Second * 5,
					BackoffCoefficient: 2.0,
					MaximumInterval:    time.Minute * 5,
					MaximumAttempts:    5,
				},
			}
			aggCtx := workflow.WithActivityOptions(ctx, aoAiAct)
			var aiAggResp *ChatCompletionQueryResponse
			// TODO, should add token chunking check here
			switch aggInst.AnalysisResponseFormat {
			case jsonFormat:
				selectTaskCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				var fullTaskDef []artemis_orchestrations.AITaskLibrary
				err = workflow.ExecuteActivity(selectTaskCtx, z.SelectTaskDefinition, ou, aggInst.AggTaskID).Get(selectTaskCtx, &fullTaskDef)
				if err != nil {
					logger.Error("failed to run agg json task selection", "Error", err)
					return err
				}
				var jdef []artemis_orchestrations.JsonSchemaDefinition
				for _, taskDef := range fullTaskDef {
					jdef = append(jdef, taskDef.Schemas...)
				}
				fname := "aggtask"
				if aggInst.AggTaskName != nil {
					fname = *aggInst.AggTaskName
				}
				if aggInst.AggModel == nil {
					logger.Error("failed to run agg json task selection", "Error", err)
					return fmt.Errorf("agg model is nil")
				}
				fd := artemis_orchestrations.ConvertToJsonDef(fname, jdef)
				jsonTaskCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				params := hera_openai.OpenAIParams{
					Model:              *aggInst.AggModel,
					FunctionDefinition: fd,
				}
				err = workflow.ExecuteActivity(jsonTaskCtx, z.CreateJsonOutputModelResponse, ou, params).Get(jsonTaskCtx, &aiAggResp)
				if err != nil {
					logger.Error("failed to run agg", "Error", err)
					return err
				}
			case socialMediaEngagementResponseFormat:
				var sg *hera_search.SearchResultGroup
				if cp.AnalysisEvalActionParams != nil && cp.AnalysisEvalActionParams.SearchResultGroup != nil {
					sg = cp.AnalysisEvalActionParams.SearchResultGroup
				}
				if sg == nil || len(sg.SearchResults) == 0 {
					continue
				}
				sg.ExtractionPromptExt = *aggInst.AggPrompt
				childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
					WorkflowID:               oj.OrchestrationName + "-agg-social-media-extraction-" + strconv.Itoa(i),
					WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
				}
				childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
				err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.SocialMediaMessagingWorkflow, ou, sg).Get(childAnalysisCtx, &aiAggResp)
				if err != nil {
					logger.Error("failed to execute child social media extraction workflow", "Error", err)
					return err
				}
			default:
				err = workflow.ExecuteActivity(aggCtx, z.AiAggregateTask, ou, aggInst, dataIn).Get(aggCtx, &aiAggResp)
				if err != nil {
					logger.Error("failed to run aggregation", "Error", err)
					return err
				}
			}
			if aiAggResp == nil || len(aiAggResp.Response.Choices) == 0 {
				continue
			}
			var aggRespId int
			aggCompCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(aggCompCtx, z.RecordCompletionResponse, ou, aiAggResp).Get(aggCompCtx, &aggRespId)
			if err != nil {
				logger.Error("failed to save agg response", "Error", err)
				return err
			}
			wr := artemis_orchestrations.AIWorkflowAnalysisResult{
				OrchestrationsID:      oj.OrchestrationID,
				ResponseID:            aggRespId,
				SourceTaskID:          *aggInst.AggTaskID,
				RunningCycleNumber:    i,
				SearchWindowUnixStart: window.UnixStartTime,
				SearchWindowUnixEnd:   window.UnixEndTime,
			}
			recordAggCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(recordAggCtx, z.SaveTaskOutput, wr, dataIn).Get(recordAggCtx, nil)
			if err != nil {
				logger.Error("failed to save aggregation resp", "Error", err)
				return err
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
			cp.WorkflowResult = wr
			ea := &EvalActionParams{
				WorkflowTemplateData: aggInst,
				ParentOutputToEval:   aiAggResp,
				EvalFns:              aggInst.AggEvalFns,
			}

			for _, evalFn := range ea.EvalFns {
				evalAggCycle := wfExecParams.CycleCountTaskRelative.AggEvalNormalizedCycleCounts[*aggInst.AggTaskID][evalFn.EvalID]
				if i%evalAggCycle == 0 {
					childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
					err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RunAiWorkflowAutoEvalProcess, cp, ea).Get(childAnalysisCtx, nil)
					if err != nil {
						logger.Error("failed to execute child analysis workflow", "Error", err)
						return err
					}
				}
			}
		}
	}
	return nil
}
