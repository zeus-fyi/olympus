package ai_platform_service_orchestrations

import (
	"strconv"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (z *ZeusAiPlatformServiceWorkflows) RunAiWorkflowChildAnalysisProcess(ctx workflow.Context, cp *MbChildSubProcessParams) error {
	if cp == nil || cp.WfExecParams.WorkflowTasks == nil || cp.Oj.OrchestrationID == 0 || cp.Ou.OrgID == 0 || cp.Ou.UserID == 0 {
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
	for _, analysisInst := range wfExecParams.WorkflowTasks {
		if runCycle%analysisInst.AnalysisCycleCount == 0 {
			if md.AnalysisRetrievals[analysisInst.AnalysisTaskID] == nil {
				continue
			}
			var sr []hera_search.SearchResult
			window := artemis_orchestrations.CalculateTimeWindowFromCycles(wfExecParams.WorkflowExecTimekeepingParams.RunWindow.UnixStartTime,
				i-analysisInst.AnalysisCycleCount, i, wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize)
			if analysisInst.RetrievalID != nil && *analysisInst.RetrievalID > 0 {
				if md.AnalysisRetrievals[analysisInst.AnalysisTaskID][*analysisInst.RetrievalID] == false {
					continue
				}
				retrievalCtx := workflow.WithActivityOptions(ctx, ao)
				err := workflow.ExecuteActivity(retrievalCtx, z.AiRetrievalTask, ou, analysisInst, window).Get(retrievalCtx, &sr)
				if err != nil {
					logger.Error("failed to run retrieval", "Error", err)
					return err
				}
				var routes []iris_models.RouteInfo
				retrievalWebCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, ou, analysisInst).Get(retrievalWebCtx, &routes)
				if err != nil {
					logger.Error("failed to run retrieval", "Error", err)
					return err
				}
				for _, route := range routes {
					fetchedResult := &hera_search.SearchResult{}
					retrievalWebTaskCtx := workflow.WithActivityOptions(ctx, aoAiAct)
					err = workflow.ExecuteActivity(retrievalWebTaskCtx, z.AiWebRetrievalTask, ou, analysisInst, route).Get(retrievalWebTaskCtx, &fetchedResult)
					if err != nil {
						logger.Error("failed to run retrieval", "Error", err)
						return err
					}
					if fetchedResult != nil && len(fetchedResult.Value) > 0 {
						sr = append(sr, *fetchedResult)
					}
				}
				md.AnalysisRetrievals[analysisInst.AnalysisTaskID][*analysisInst.RetrievalID] = false
				if len(sr) == 0 {
					continue
				}
			}
			analysisCtx := workflow.WithActivityOptions(ctx, aoAiAct)
			var aiResp *ChatCompletionQueryResponse
			err := workflow.ExecuteActivity(analysisCtx, z.AiAnalysisTask, ou, analysisInst, sr).Get(analysisCtx, &aiResp)
			if err != nil {
				logger.Error("failed to run analysis", "Error", err)
				return err
			}
			if aiResp == nil || len(aiResp.Response.Choices) == 0 {
				continue
			}
			var analysisRespId int
			analysisCompCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(analysisCompCtx, z.RecordCompletionResponse, ou, aiResp).Get(analysisCompCtx, &analysisRespId)
			if err != nil {
				logger.Error("failed to save analysis response", "Error", err)
				return err
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
				return err
			}
			evalWfID := oj.OrchestrationName + "-analysis-eval=" + strconv.Itoa(i)
			childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
				WorkflowID:               oj.OrchestrationName + "-analysis-eval=" + strconv.Itoa(i),
				WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
			}
			cp.Window = window
			cp.WfID = evalWfID
			cp.WorkflowResult = wr

			ea := &EvalActionParams{
				WorkflowTemplateData: analysisInst,
				ParentOutputToEval:   aiResp,
				EvalFns:              analysisInst.AnalysisTaskDB.AnalysisEvalFns,
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
						return err
					}
				}
			}
		}
	}
	return nil
}
