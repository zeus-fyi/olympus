package ai_platform_service_orchestrations

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (z *ZeusAiPlatformServiceWorkflows) RunAiChildAnalysisProcessWorkflow(ctx workflow.Context, cp *MbChildSubProcessParams) error {
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
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    25,
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
				childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
					WorkflowID:               CreateExecAiWfId(oj.OrchestrationName + "-analysis-ret-" + strconv.Itoa(i)),
					WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
					RetryPolicy:              ao.RetryPolicy,
				}
				var rets []artemis_orchestrations.RetrievalItem
				chunkedTaskCtx := workflow.WithActivityOptions(ctx, ao)
				err := workflow.ExecuteActivity(chunkedTaskCtx, z.SelectRetrievalTask, ou, *analysisInst.RetrievalID).Get(chunkedTaskCtx, &rets)
				if err != nil {
					logger.Error("failed to run analysis json", "Error", err)
					return err
				}
				if len(rets) <= 0 {
					continue
				}
				tte := TaskToExecute{
					WfID: childAnalysisWorkflowOptions.WorkflowID,
					Ou:   ou,
					Wft:  analysisInst,
					Sg:   sg,
					Tc: TaskContext{
						Model:     analysisInst.AnalysisModel,
						TaskID:    analysisInst.AnalysisTaskID,
						Retrieval: rets[0],
					},
				}
				childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
				err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RetrievalsWorkflow, tte).Get(childAnalysisCtx, &sg)
				if err != nil {
					logger.Error("failed to execute child retrieval workflow", "Error", err)
					return err
				}
				md.AnalysisRetrievals[analysisInst.AnalysisTaskID][*analysisInst.RetrievalID] = false
				if len(sg.SearchResults) == 0 && len(sg.ApiResponseResults) <= 0 {
					continue
				}
			}
			pr := &PromptReduction{TokenOverflowStrategy: analysisInst.AnalysisTokenOverflowStrategy, Model: analysisInst.AnalysisModel, MarginBuffer: analysisInst.AnalysisMarginBuffer}
			if sg != nil && len(sg.SearchResults) > 0 {
				pr.PromptReductionSearchResults = &PromptReductionSearchResults{
					InPromptBody:  analysisInst.AnalysisPrompt,
					InSearchGroup: sg,
				}
			} else if sg != nil && len(sg.ApiResponseResults) > 0 {
				pr.PromptReductionSearchResults = &PromptReductionSearchResults{
					InPromptBody:  analysisInst.AnalysisPrompt,
					InSearchGroup: sg,
				}
			} else {
				pr.PromptReductionText = &PromptReductionText{
					InPromptBody: analysisInst.AnalysisPrompt,
				}
			}
			chunkedTaskCtx := workflow.WithActivityOptions(ctx, ao)
			err := workflow.ExecuteActivity(chunkedTaskCtx, z.TokenOverflowReduction, ou, pr).Get(chunkedTaskCtx, &pr)
			if err != nil {
				logger.Error("failed to run analysis json", "Error", err)
				return err
			}

			chunkIterator := getChunkIteratorLen(pr)
			for chunkOffset := 0; chunkOffset < chunkIterator; chunkOffset++ {
				if pr.PromptReductionSearchResults != nil && pr.PromptReductionSearchResults.OutSearchGroups != nil && chunkOffset < len(pr.PromptReductionSearchResults.OutSearchGroups) {
					sg = pr.PromptReductionSearchResults.OutSearchGroups[chunkOffset]
					sg.Model = analysisInst.AnalysisModel
					sg.ResponseFormat = analysisInst.AnalysisResponseFormat
				} else {
					sg = &hera_search.SearchResultGroup{
						Model:          analysisInst.AnalysisModel,
						ResponseFormat: analysisInst.AnalysisResponseFormat,
						BodyPrompt:     pr.PromptReductionText.OutPromptChunks[chunkOffset],
						SearchResults:  []hera_search.SearchResult{},
					}
				}
				wr := &artemis_orchestrations.AIWorkflowAnalysisResult{OrchestrationID: oj.OrchestrationID, SourceTaskID: analysisInst.AnalysisTaskID, RunningCycleNumber: i, ChunkOffset: chunkOffset, IterationCount: 0, SearchWindowUnixStart: window.UnixStartTime, SearchWindowUnixEnd: window.UnixEndTime}
				tte := TaskToExecute{Ou: ou, Wft: analysisInst, Sg: sg, Wr: wr}
				var analysisRespId int
				var aiResp *ChatCompletionQueryResponse
				switch analysisInst.AnalysisResponseFormat {
				case jsonFormat, socialMediaExtractionResponseFormat:
					tte.Sg.ExtractionPromptExt = analysisInst.AnalysisPrompt
					tte.Sg.SourceTaskID = analysisInst.AnalysisTaskID
					childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
						WorkflowID:               CreateExecAiWfId(oj.OrchestrationName + fmt.Sprintf("-analysis-%s-task-%d-chunk-%d", analysisInst.AnalysisResponseFormat, i, chunkOffset)),
						WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
						RetryPolicy:              ao.RetryPolicy,
					}
					tte.WfID = childAnalysisWorkflowOptions.WorkflowID
					tte.Tc = TaskContext{TaskName: analysisInst.AnalysisTaskName, TaskType: AnalysisTask, ResponseFormat: analysisInst.AnalysisResponseFormat, Model: analysisInst.AnalysisModel, TaskID: analysisInst.AnalysisTaskID}
					var fullTaskDef []artemis_orchestrations.AITaskLibrary
					selectTaskCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(selectTaskCtx, z.SelectTaskDefinition, tte.Ou, tte.Sg.SourceTaskID).Get(selectTaskCtx, &fullTaskDef)
					if err != nil {
						logger.Error("failed to run task", "Error", err)
						return err
					}
					if len(fullTaskDef) == 0 {
						logger.Warn("failed to run task", "Error", err)
						continue
					}
					var jdef []*artemis_orchestrations.JsonSchemaDefinition
					for _, taskDef := range fullTaskDef {
						jdef = append(jdef, taskDef.Schemas...)
					}
					tte.Tc.Schemas = jdef
					childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
					err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.JsonOutputTaskWorkflow, tte).Get(childAnalysisCtx, &aiResp)
					if err != nil {
						logger.Error("failed to execute analysis json workflow", "Error", err)
						return err
					}
					if aiResp != nil {
						if cp.AnalysisEvalActionParams == nil {
							cp.AnalysisEvalActionParams = &EvalActionParams{}
						}
						cp.AnalysisEvalActionParams.ParentOutputToEval = aiResp
						cp.AnalysisEvalActionParams.SearchResultGroup = tte.Sg
					} else {
						continue
					}
				default:
					inGroup := tte.Sg.SearchResults
					if len(tte.Sg.ApiResponseResults) > 0 {
						inGroup = tte.Sg.ApiResponseResults
					}
					tte.Tc = TaskContext{TaskName: analysisInst.AnalysisTaskName, TaskType: AnalysisTask, Temperature: float32(analysisInst.AnalysisTemperature), ResponseFormat: analysisInst.AnalysisResponseFormat, Model: analysisInst.AnalysisModel, TaskID: analysisInst.AnalysisTaskID}
					analysisCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(analysisCtx, z.AiAnalysisTask, ou, analysisInst, inGroup).Get(analysisCtx, &aiResp)
					if err != nil {
						logger.Error("failed to run analysis", "Error", err)
						return err
					}
					analysisCompCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(analysisCompCtx, z.RecordCompletionResponse, ou, aiResp).Get(analysisCompCtx, &analysisRespId)
					if err != nil {
						logger.Error("failed to save analysis response", "Error", err)
						return err
					}
					wr.ResponseID = analysisRespId
					afv := InputDataAnalysisToAgg{
						ChatCompletionQueryResponse: aiResp,
						SearchResultGroup:           tte.Sg,
					}
					recordAnalysisCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(recordAnalysisCtx, z.SaveTaskOutput, &wr, afv).Get(recordAnalysisCtx, &aiResp.WorkflowResultID)
					if err != nil {
						logger.Error("failed to save analysis", "Error", err)
						return err
					}
					if aiResp != nil && aiResp.Prompt != nil {
						tte.Sg.ResponseBody = aiResp.Prompt["response"]
					}
					tte.Sg.BodyPrompt = hera_search.FormatSearchResultsV2(inGroup)
				}
				if aiResp == nil || len(aiResp.Response.Choices) == 0 {
					continue
				}
				// TODO, run in parallel
				for ind, evalFn := range analysisInst.AnalysisTaskDB.AnalysisEvalFns {
					if evalFn.EvalID == 0 {
						continue
					}
					var evalAnalysisOnlyCycle int
					if analysisInst.AggTaskID != nil {
						evalAnalysisOnlyCycle = wfExecParams.CycleCountTaskRelative.AggAnalysisEvalNormalizedCycleCounts[*analysisInst.AggTaskID][analysisInst.AnalysisTaskID][evalFn.EvalID]
					} else {
						evalAnalysisOnlyCycle = wfExecParams.CycleCountTaskRelative.AnalysisEvalNormalizedCycleCounts[analysisInst.AnalysisTaskID][evalFn.EvalID]
					}
					if i%evalAnalysisOnlyCycle == 0 {
						childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
							WorkflowID:               CreateExecAiWfId(oj.OrchestrationName + "-analysis-eval-" + strconv.Itoa(i) + "-chunk-" + strconv.Itoa(chunkOffset) + "eval-fn" + strconv.Itoa(evalFn.EvalID) + "-ind-" + strconv.Itoa(ind)),
							WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
							RetryPolicy:              ao.RetryPolicy,
						}
						cp.Window = window
						cp.WfID = childAnalysisWorkflowOptions.WorkflowID
						if wr.WorkflowResultID == 0 {
							wr.WorkflowResultID = aiResp.WorkflowResultID
						}
						cp.WorkflowResult = *wr
						ea := &EvalActionParams{WorkflowTemplateData: analysisInst, ParentOutputToEval: aiResp, EvalFns: analysisInst.AnalysisTaskDB.AnalysisEvalFns, SearchResultGroup: tte.Sg, TaskToExecute: tte}
						if analysisInst.AggTaskID != nil && analysisInst.AggAnalysisEvalFns != nil {
							ea.EvalFns = analysisInst.AggAnalysisEvalFns
						} else if analysisInst.AnalysisEvalFns != nil {
							ea.EvalFns = analysisInst.AnalysisTaskDB.AnalysisEvalFns
						}
						childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
						wio := &WorkflowStageIO{
							WorkflowStageReference: artemis_orchestrations.WorkflowStageReference{
								WorkflowRunID: oj.OrchestrationID,
								ChildWfID:     childAnalysisWorkflowOptions.WorkflowID,
								RunCycle:      runCycle,
							},
							WorkflowStageInfo: WorkflowStageInfo{
								RunAiWorkflowAutoEvalProcessInputs: &RunAiWorkflowAutoEvalProcessInputs{
									Mb:  cp,
									Cpe: ea,
								},
							},
						}
						saveWfStageIOCtx := workflow.WithActivityOptions(ctx, ao)
						err = workflow.ExecuteActivity(saveWfStageIOCtx, z.SaveWorkflowIO, wio).Get(saveWfStageIOCtx, &wio)
						if err != nil {
							logger.Error("failed to saveWfStageIOCtx results", "Error", err)
							return err
						}
						if wio == nil || wio.InputID == 0 {
							err = fmt.Errorf("wio.InputID is 0 in analysis eval")
							logger.Warn("wio.InputID is 0 in analysis eval", "Error", err)
							return err
						}
						err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RunAiWorkflowAutoEvalProcess, wio.InputID).Get(childAnalysisCtx, nil)
						if err != nil {
							logger.Error("failed to execute child analysis workflow", "Error", err)
							return err
						}
					}
				}
			}
		}
	}
	return nil
}
