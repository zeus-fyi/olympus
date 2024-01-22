package ai_platform_service_orchestrations

import (
	"strconv"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
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
	// TODO update activity options by wfExecParams
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
				retWfID := oj.OrchestrationName + "-analysis-ret-" + strconv.Itoa(i)
				childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
					WorkflowID:               retWfID,
					WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
				}
				tte := TaskToExecute{
					WfID: retWfID,
					Ou:   ou,
					Wft:  analysisInst,
					Sg:   sg,
				}
				childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
				err := workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RetrievalsWorkflow, cp, tte).Get(childAnalysisCtx, &sg)
				if err != nil {
					logger.Error("failed to execute child retrieval workflow", "Error", err)
					return nil, err
				}
				md.AnalysisRetrievals[analysisInst.AnalysisTaskID][*analysisInst.RetrievalID] = false
				if len(sg.SearchResults) == 0 {
					continue
				}
			}
			pr := &PromptReduction{
				TokenOverflowStrategy: analysisInst.AnalysisTokenOverflowStrategy,
				Model:                 analysisInst.AnalysisModel,
			}
			if sg != nil && len(sg.SearchResults) > 0 {
				pr.PromptReductionSearchResults = &PromptReductionSearchResults{
					InPromptBody:  analysisInst.AnalysisPrompt,
					InSearchGroup: sg,
				}
			} else {
				pr.PromptReductionText = &PromptReductionText{
					InPromptBody: analysisInst.AnalysisPrompt,
				}
			}
			chunkedTaskCtx := workflow.WithActivityOptions(ctx, aoAiAct)
			err := workflow.ExecuteActivity(chunkedTaskCtx, z.TokenOverflowReduction, ou, pr).Get(chunkedTaskCtx, &pr)
			if err != nil {
				logger.Error("failed to run analysis json", "Error", err)
				return nil, err
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
					analysisInst.AnalysisPrompt = pr.PromptReductionText.OutPromptChunks[chunkOffset]
				}

				wr := &artemis_orchestrations.AIWorkflowAnalysisResult{
					OrchestrationsID:      oj.OrchestrationID,
					SourceTaskID:          analysisInst.AnalysisTaskID,
					RunningCycleNumber:    i,
					ChunkOffset:           chunkOffset,
					IterationCount:        0,
					SearchWindowUnixStart: window.UnixStartTime,
					SearchWindowUnixEnd:   window.UnixEndTime,
				}

				var aiResp *ChatCompletionQueryResponse
				tte := TaskToExecute{
					Ou:  ou,
					Wft: analysisInst,
					Sg:  sg,
					Wr:  wr,
				}
				var wfID string
				var analysisRespId int
				switch analysisInst.AnalysisResponseFormat {
				case jsonFormat:
					sg.ExtractionPromptExt = analysisInst.AnalysisPrompt
					sg.SourceTaskID = analysisInst.AnalysisTaskID
					wfID = oj.OrchestrationName + "-analysis-json-task-" + strconv.Itoa(i)
					tte.WfID = wfID
					pr.Model = tte.Wft.AnalysisModel
					pr.TokenOverflowStrategy = tte.Wft.AnalysisTokenOverflowStrategy
					pr.PromptReductionText.InPromptBody = tte.Wft.AnalysisPrompt
					childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
						WorkflowID:               wfID,
						WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
					}
					tte.Tc = TaskContext{
						TaskName:    analysisInst.AnalysisTaskName,
						TaskType:    AnalysisTask,
						Model:       analysisInst.AnalysisModel,
						TaskID:      analysisInst.AnalysisTaskID,
						ChunkOffset: chunkOffset,
					}
					childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
					err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.JsonOutputTaskWorkflow, tte).Get(childAnalysisCtx, &aiResp)
					if err != nil {
						logger.Error("failed to execute analysis json workflow", "Error", err)
						return nil, err
					}
					sg.SearchResults = aiResp.FilteredSearchResults
				case socialMediaExtractionResponseFormat:
					if sg == nil || len(sg.SearchResults) == 0 {
						continue
					}
					sg.ExtractionPromptExt = analysisInst.AnalysisPrompt
					wfID = oj.OrchestrationName + "-analysis-social-media-extraction-" + strconv.Itoa(i)
					childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
						WorkflowID:               wfID,
						WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
					}
					childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
					err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.SocialMediaExtractionWorkflow, tte).Get(childAnalysisCtx, &aiResp)
					if err != nil {
						logger.Error("failed to execute child social media extraction workflow", "Error", err)
						return nil, err
					}
					sg.SearchResults = aiResp.FilteredSearchResults
					analysisCompCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(analysisCompCtx, z.RecordCompletionResponse, ou, aiResp).Get(analysisCompCtx, &analysisRespId)
					if err != nil {
						logger.Error("failed to save analysis response", "Error", err)
						return nil, err
					}
					recordAnalysisCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(recordAnalysisCtx, z.SaveTaskOutput, &wr).Get(recordAnalysisCtx, nil)
					if err != nil {
						logger.Error("failed to save analysis", "Error", err)
						return nil, err
					}
				default:
					analysisCtx := workflow.WithActivityOptions(ctx, aoAiAct)
					err = workflow.ExecuteActivity(analysisCtx, z.AiAnalysisTask, ou, analysisInst, sg.SearchResults).Get(analysisCtx, &aiResp)
					if err != nil {
						logger.Error("failed to run analysis", "Error", err)
						return nil, err
					}
					analysisCompCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(analysisCompCtx, z.RecordCompletionResponse, ou, aiResp).Get(analysisCompCtx, &analysisRespId)
					if err != nil {
						logger.Error("failed to save analysis response", "Error", err)
						return nil, err
					}
					recordAnalysisCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(recordAnalysisCtx, z.SaveTaskOutput, &wr).Get(recordAnalysisCtx, nil)
					if err != nil {
						logger.Error("failed to save analysis", "Error", err)
						return nil, err
					}
				}
				if aiResp == nil || len(aiResp.Response.Choices) == 0 {
					continue
				}
				wfID = oj.OrchestrationName + "-analysis-eval-" + strconv.Itoa(i)
				childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
					WorkflowID:               wfID,
					WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
				}
				cp.Window = window
				cp.WfID = wfID
				cp.WorkflowResult = *wr
				ea := &EvalActionParams{
					WorkflowTemplateData: analysisInst,
					ParentOutputToEval:   aiResp,
					EvalFns:              analysisInst.AnalysisTaskDB.AnalysisEvalFns,
					SearchResultGroup:    sg,
					TaskToExecute:        tte,
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
	}
	return cp, nil
}
