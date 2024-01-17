package ai_platform_service_orchestrations

import (
	"strconv"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (z *ZeusAiPlatformServiceWorkflows) RunAiWorkflowChildAnalysisProcess(ctx workflow.Context, cp *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	if cp == nil || cp.WfExecParams.WorkflowTasks == nil || cp.Oj.OrchestrationID == 0 || cp.Ou.OrgID == 0 || cp.Ou.UserID == 0 {
		return nil, nil
	}
	wfExecParams := cp.WfExecParams
	ou := cp.Ou
	oj := cp.Oj
	runCycle := cp.RunCycle

	md := artemis_orchestrations.MapDependencies(wfExecParams.WorkflowTasks)
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15, // Setting a valid non-zero timeout
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
	i := runCycle
	sg := &hera_search.SearchResultGroup{
		SearchResults: []hera_search.SearchResult{},
	}
	for _, analysisInst := range wfExecParams.WorkflowTasks {
		if runCycle%analysisInst.AnalysisCycleCount == 0 {
			if md.AnalysisRetrievals[analysisInst.AnalysisTaskID] == nil {
				continue
			}
			window := artemis_orchestrations.CalculateTimeWindowFromCycles(wfExecParams.WorkflowExecTimekeepingParams.RunWindow.UnixStartTime,
				i-analysisInst.AnalysisCycleCount, i, wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize)
			if analysisInst.RetrievalID != nil && *analysisInst.RetrievalID > 0 {
				if md.AnalysisRetrievals[analysisInst.AnalysisTaskID][*analysisInst.RetrievalID] == false {
					continue
				}
				sg = &hera_search.SearchResultGroup{
					PlatformName:   analysisInst.RetrievalPlatform,
					Model:          analysisInst.AnalysisModel,
					ResponseFormat: analysisInst.AnalysisResponseFormat,
					SearchResults:  []hera_search.SearchResult{},
					Window:         window,
				}
				switch analysisInst.RetrievalPlatform {
				case twitterPlatform, redditPlatform, discordPlatform, telegramPlatform:
					retrievalCtx := workflow.WithActivityOptions(ctx, ao)
					var sr []hera_search.SearchResult
					err := workflow.ExecuteActivity(retrievalCtx, z.AiRetrievalTask, ou, analysisInst, window).Get(retrievalCtx, &sr)
					if err != nil {
						logger.Error("failed to run retrieval", "Error", err)
						return nil, err
					}
					sg.SearchResults = append(sg.SearchResults, sr...)
				case webPlatform:
					var routes []iris_models.RouteInfo
					retrievalWebCtx := workflow.WithActivityOptions(ctx, aoAiAct)
					err := workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, ou, analysisInst).Get(retrievalWebCtx, &routes)
					if err != nil {
						logger.Error("failed to run retrieval", "Error", err)
						return nil, err
					}
					for _, route := range routes {
						fetchedResult := &hera_search.SearchResult{}
						retrievalWebTaskCtx := workflow.WithActivityOptions(ctx, aoAiAct)
						err = workflow.ExecuteActivity(retrievalWebTaskCtx, z.AiWebRetrievalTask, ou, analysisInst, route).Get(retrievalWebTaskCtx, &fetchedResult)
						if err != nil {
							logger.Error("failed to run retrieval", "Error", err)
							return nil, err
						}
						if fetchedResult != nil && len(fetchedResult.Value) > 0 {
							sg.SearchResults = append(sg.SearchResults, *fetchedResult)
						}
					}
				}
				md.AnalysisRetrievals[analysisInst.AnalysisTaskID][*analysisInst.RetrievalID] = false
				if len(sg.SearchResults) == 0 {
					continue
				}
			}

			var aiResp *ChatCompletionQueryResponse
			var fullTaskDef []artemis_orchestrations.AITaskLibrary
			if analysisInst.AnalysisResponseFormat == socialMediaExtractionResponseFormat || analysisInst.AnalysisResponseFormat == jsonFormat {
				selectTaskCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				err := workflow.ExecuteActivity(selectTaskCtx, z.SelectTaskDefinition, ou, analysisInst.AnalysisTaskID).Get(selectTaskCtx, &fullTaskDef)
				if err != nil {
					logger.Error("failed to run analysis", "Error", err)
					return nil, err
				}
				if len(fullTaskDef) == 0 {
					continue
				}
			}
			// TODO, should add token chunking check here
			switch analysisInst.AnalysisResponseFormat {
			case jsonFormat:
				var jdef []artemis_orchestrations.JsonSchemaDefinition
				for _, taskDef := range fullTaskDef {
					jdef = append(jdef, taskDef.Schemas...)
				}
				fd := artemis_orchestrations.ConvertToJsonDef(analysisInst.AnalysisTaskName, jdef)
				jsonTaskCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				params := hera_openai.OpenAIParams{
					Model:              analysisInst.AnalysisModel,
					FunctionDefinition: fd,
				}
				err := workflow.ExecuteActivity(jsonTaskCtx, z.CreateJsonOutputModelResponse, ou, params).Get(jsonTaskCtx, &aiResp)
				if err != nil {
					logger.Error("failed to run analysis json", "Error", err)
					return nil, err
				}
			case socialMediaExtractionResponseFormat:
				if sg == nil || len(sg.SearchResults) == 0 {
					continue
				}
				sg.ExtractionPromptExt = analysisInst.AnalysisPrompt
				childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
					WorkflowID:               oj.OrchestrationName + "-analysis-social-media-extraction-" + strconv.Itoa(i),
					WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
				}
				childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
				err := workflow.ExecuteChildWorkflow(childAnalysisCtx, z.SocialMediaExtractionWorkflow, ou, sg).Get(childAnalysisCtx, &aiResp)
				if err != nil {
					logger.Error("failed to execute child social media extraction workflow", "Error", err)
					return nil, err
				}
				// TODO, now these extraction results are going to eval, check process
			default:
				analysisCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				err := workflow.ExecuteActivity(analysisCtx, z.AiAnalysisTask, ou, analysisInst, sg.SearchResults).Get(analysisCtx, &aiResp)
				if err != nil {
					logger.Error("failed to run analysis", "Error", err)
					return nil, err
				}
			}
			if aiResp == nil || len(aiResp.Response.Choices) == 0 {
				continue
			}
			// TODO, validate chunk saving works
			var analysisRespId int
			analysisCompCtx := workflow.WithActivityOptions(ctx, ao)
			err := workflow.ExecuteActivity(analysisCompCtx, z.RecordCompletionResponse, ou, aiResp).Get(analysisCompCtx, &analysisRespId)
			if err != nil {
				logger.Error("failed to save analysis response", "Error", err)
				return nil, err
			}
			wr := artemis_orchestrations.AIWorkflowAnalysisResult{
				OrchestrationsID:      oj.OrchestrationID,
				ResponseID:            analysisRespId,
				SourceTaskID:          analysisInst.AnalysisTaskID,
				RunningCycleNumber:    i,
				SearchWindowUnixStart: window.UnixStartTime,
				SearchWindowUnixEnd:   window.UnixEndTime,
			}
			recordAnalysisCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(recordAnalysisCtx, z.SaveTaskOutput, wr).Get(recordAnalysisCtx, nil)
			if err != nil {
				logger.Error("failed to save analysis", "Error", err)
				return nil, err
			}
			evalWfID := oj.OrchestrationName + "-analysis-eval-" + strconv.Itoa(i)
			childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
				WorkflowID:               oj.OrchestrationName + "-analysis-eval-" + strconv.Itoa(i),
				WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
			}
			cp.Window = window
			cp.WfID = evalWfID
			cp.WorkflowResult = wr

			ea := &EvalActionParams{
				WorkflowTemplateData: analysisInst,
				ParentOutputToEval:   aiResp,
				EvalFns:              analysisInst.AnalysisTaskDB.AnalysisEvalFns,
				SearchResultGroup:    sg,
			}
			if analysisInst.AggTaskID != nil {
				ea.EvalFns = analysisInst.AggAnalysisEvalFns
			}
			for _, evalFn := range ea.EvalFns {
				var evalAnalysisOnlyCycle int
				if analysisInst.AggTaskID != nil {
					evalAnalysisOnlyCycle = wfExecParams.CycleCountTaskRelative.AggAnalysisEvalNormalizedCycleCounts[*analysisInst.AggTaskID][analysisInst.AnalysisTaskID][evalFn.EvalID]
				} else {
					evalAnalysisOnlyCycle = wfExecParams.CycleCountTaskRelative.AnalysisEvalNormalizedCycleCounts[analysisInst.AnalysisTaskID][evalFn.EvalID]
				}
				if i%evalAnalysisOnlyCycle == 0 {
					childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
					err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RunAiWorkflowAutoEvalProcess, cp, ea).Get(childAnalysisCtx, nil)
					if err != nil {
						logger.Error("failed to execute child analysis workflow", "Error", err)
						return nil, err
					}
				}
			}
			cp.AnalysisEvalActionParams = ea
			return cp, nil
		}
	}
	return cp, nil
}
