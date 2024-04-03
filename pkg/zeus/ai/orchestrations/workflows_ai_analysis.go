package ai_platform_service_orchestrations

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (z *ZeusAiPlatformServiceWorkflows) RunAiChildAnalysisProcessWorkflow(ctx workflow.Context, cp *MbChildSubProcessParams) error {
	if cp == nil || cp.WfExecParams.WorkflowTasks == nil || cp.Oj.OrchestrationID == 0 || cp.Ou.OrgID == 0 || cp.Ou.UserID == 0 {
		return nil
	}
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Hour * 24, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 3,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    100,
		},
	}
	wfExecParams := cp.WfExecParams
	ou := cp.Ou
	oj := cp.Oj
	runCycle := cp.Wsr.RunCycle
	md := artemis_orchestrations.MapDependencies(wfExecParams.WorkflowTasks)
	i := runCycle
	log.Info().Interface("runCycle", runCycle).Msg("analysis: runCycle")
	for _, analysisInst := range wfExecParams.WorkflowTasks {
		// if agg task id; meant only for agg tasks; since agg task can contain original related analysis task embedded
		if analysisInst.AggTaskID != nil {
			continue
		}
		// checks for run cycle validity
		if runCycle%analysisInst.AnalysisCycleCount == 0 {
			if md.AnalysisRetrievals[analysisInst.AnalysisTaskID] == nil {
				continue
			}
			log.Info().Interface("taskID", analysisInst.AnalysisTaskID).Interface("taskName", cp.Tc.TaskName).Msg("analysis: taskID")
			window := artemis_orchestrations.CalculateTimeWindowFromCycles(wfExecParams.WorkflowExecTimekeepingParams.RunWindow.UnixStartTime,
				i-analysisInst.AnalysisCycleCount, i, wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize)
			cp.Window = window
			cp.Tc = getAnalysisTaskContext(analysisInst)
			// todo revist placement
			wsrCreateCtx := workflow.WithActivityOptions(ctx, ao)
			err := workflow.ExecuteActivity(wsrCreateCtx, z.CreateWsr, cp).Get(wsrCreateCtx, &cp)
			if err != nil {
				logger.Error("failed to run get retrieval routes", "Error", err)
				return err
			}
			if analysisInst.RetrievalID != nil && *analysisInst.RetrievalID > 0 {
				if md.AnalysisRetrievals[analysisInst.AnalysisTaskID][*analysisInst.RetrievalID] == false {
					continue
				}
				var rets []artemis_orchestrations.RetrievalItem
				tmpOu := ou
				if wfExecParams.WorkflowOverrides.IsUsingFlows {
					tmpOu.OrgID = FlowsOrgID
				}
				chunkedTaskCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(chunkedTaskCtx, z.SelectRetrievalTask, tmpOu, *analysisInst.RetrievalID).Get(chunkedTaskCtx, &rets)
				if err != nil {
					logger.Error("failed to run analysis retrieval", "Error", err)
					return err
				}
				if len(rets) <= 0 {
					continue
				}
				cp.Tc.Retrieval = rets[0]
				aoRet := workflow.ActivityOptions{
					ScheduleToCloseTimeout: time.Hour * 24, // Setting a valid non-zero timeout
					RetryPolicy: &temporal.RetryPolicy{
						InitialInterval:    time.Second * 3,
						BackoffCoefficient: 2.0,
						MaximumInterval:    time.Minute * 5,
						MaximumAttempts:    10000,
					},
				}
				childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
					WorkflowID:         oj.OrchestrationName + "-analysis-ret-cycle-" + strconv.Itoa(runCycle),
					WorkflowRunTimeout: ao.ScheduleToCloseTimeout,
					RetryPolicy:        aoRet.RetryPolicy,
				}
				cp.Wsr.ChildWfID = childAnalysisWorkflowOptions.WorkflowID
				childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
				err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RetrievalsWorkflow, cp).Get(childAnalysisCtx, &cp)
				if err != nil {
					logger.Error("failed to execute child retrieval workflow", "Error", err)
					return err
				}
				md.AnalysisRetrievals[analysisInst.AnalysisTaskID][*analysisInst.RetrievalID] = false
			}
			pr := &PromptReduction{
				MarginBuffer:          analysisInst.AnalysisMarginBuffer,
				Model:                 analysisInst.AnalysisModel,
				TokenOverflowStrategy: analysisInst.AnalysisTokenOverflowStrategy,
				PromptReductionText: &PromptReductionText{
					InPromptBody: analysisInst.AnalysisPrompt,
				},
			}
			log.Info().Msg("analysis: running token overflow reduction")
			var chunkIterator int
			chunkedTaskCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(chunkedTaskCtx, z.TokenOverflowReduction, cp, pr).Get(chunkedTaskCtx, &cp)
			if err != nil {
				logger.Error("failed to run analysis token overflow", "Error", err)
				return err
			}
			chunkIterator = cp.Tc.ChunkIterator
			log.Info().Int("chunkIterator", chunkIterator).Msg("analysis: chunkIterator")
			for chunkOffset := 0; chunkOffset < chunkIterator; chunkOffset++ {
				cp.Wsr.ChunkOffset = chunkOffset
				log.Info().Interface("chunkOffset", chunkOffset).Msg("analysis: chunkOffset")
				switch analysisInst.AnalysisResponseFormat {
				case jsonFormat:
					childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
						WorkflowID:               oj.OrchestrationName + fmt.Sprintf("-analysis-%s-task-%d-chunk-%d", analysisInst.AnalysisResponseFormat, i, chunkOffset),
						WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
						RetryPolicy:              ao.RetryPolicy,
					}
					cp.Wsr.ChildWfID = childAnalysisWorkflowOptions.WorkflowID
					childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
					err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.JsonOutputTaskWorkflow, cp).Get(childAnalysisCtx, &cp)
					if err != nil {
						logger.Error("failed to execute analysis json workflow", "Error", err)
						return err
					}
				case csvFormat:
					analysisCompCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(analysisCompCtx, z.RecordCompletionResponse, ou, getDummyChatCompResp()).Get(analysisCompCtx, &cp.Tc.ResponseID)
					if err != nil {
						logger.Error("failed to save analysis read only response", "Error", err)
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
						ResponseID:            cp.Tc.ResponseID,
					}
					cp.Tc.ResponseFormat = csvFormat
					recordAnalysisCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(recordAnalysisCtx, z.SaveCsvTaskOutput, cp, wr).Get(recordAnalysisCtx, &cp.Tc.WorkflowResultID)
					if err != nil {
						logger.Error("failed to save csv analysis", "Error", err)
						return err
					}
				case sendNextStage:
				default:
					var aiResp *ChatCompletionQueryResponse
					analysisCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(analysisCtx, z.AiAnalysisTask, ou, analysisInst, cp).Get(analysisCtx, &aiResp)
					if err != nil {
						logger.Error("failed to run analysis", "Error", err)
						return err
					}
					analysisCompCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(analysisCompCtx, z.RecordCompletionResponse, ou, aiResp).Get(analysisCompCtx, &cp.Tc.ResponseID)
					if err != nil {
						logger.Error("failed to save analysis response", "Error", err)
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
						ResponseID:            cp.Tc.ResponseID,
					}
					ia := InputDataAnalysisToAgg{
						ChatCompletionQueryResponse: aiResp,
					}
					var tmp string
					for _, cv := range aiResp.Response.Choices {
						tmp += cv.Message.Content + "\n"
					}
					ia.TextInput = &tmp
					recordAnalysisCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(recordAnalysisCtx, z.SaveTaskOutput, wr, cp, ia).Get(recordAnalysisCtx, &cp.Tc.WorkflowResultID)
					if err != nil {
						logger.Error("failed to save analysis", "Error", err)
						return err
					}
				}
				logger.Info("analysis: len(evalFns)", len(analysisInst.AnalysisTaskDB.AnalysisEvalFns))
				for ind, evalFn := range analysisInst.AnalysisTaskDB.AnalysisEvalFns {
					if evalFn.EvalID == 0 {
						continue
					}
					logger.Info("analysis: eval", evalFn.EvalID)
					var evalAnalysisOnlyCycle int
					if analysisInst.AggTaskID != nil {
						evalAnalysisOnlyCycle = wfExecParams.CycleCountTaskRelative.AggAnalysisEvalNormalizedCycleCounts[*analysisInst.AggTaskID][analysisInst.AnalysisTaskID][evalFn.EvalID]
					} else {
						evalAnalysisOnlyCycle = wfExecParams.CycleCountTaskRelative.AnalysisEvalNormalizedCycleCounts[analysisInst.AnalysisTaskID][evalFn.EvalID]
					}
					if i%evalAnalysisOnlyCycle == 0 {
						childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
							WorkflowID:         oj.OrchestrationName + "-analysis-eval-" + strconv.Itoa(i) + "-chunk-" + strconv.Itoa(chunkOffset) + "-eval-fn" + strconv.Itoa(evalFn.EvalID) + "-ind-" + strconv.Itoa(ind),
							WorkflowRunTimeout: ao.ScheduleToCloseTimeout,
							RetryPolicy:        ao.RetryPolicy,
						}
						cp.Tc.EvalID = evalFn.EvalID
						log.Info().Msg("running analysis eval")
						childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
						err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RunAiWorkflowAutoEvalProcess, cp).Get(childAnalysisCtx, nil)
						if err != nil {
							logger.Error("failed to execute child analysis workflow", "Error", err)
							return err
						}
					}
				}
			}
			log.Info().Msg("analysis: evalFns complete")
		}
		log.Info().Interface("runCycle", runCycle).Msg("analysis: runCycle done")
	}
	return nil
}

func getDummyChatCompResp() *ChatCompletionQueryResponse {
	aiResp := &ChatCompletionQueryResponse{
		Prompt: make(map[string]string),
		Response: openai.ChatCompletionResponse{
			Model: "none",
			Usage: openai.Usage{
				PromptTokens:     0,
				CompletionTokens: 0,
				TotalTokens:      0,
			},
		},
	}
	return aiResp
}

func getAnalysisTaskContext(analysisInst artemis_orchestrations.WorkflowTemplateData) TaskContext {
	return TaskContext{
		TaskName:              analysisInst.AnalysisTaskName,
		TaskType:              AnalysisTask,
		Model:                 analysisInst.AnalysisModel,
		TaskID:                analysisInst.AnalysisTaskID,
		ResponseFormat:        analysisInst.AnalysisResponseFormat,
		Prompt:                analysisInst.AnalysisPrompt,
		WorkflowTemplateData:  analysisInst,
		TokenOverflowStrategy: analysisInst.AnalysisTokenOverflowStrategy,
		MarginBuffer:          analysisInst.AnalysisMarginBuffer,
		Temperature:           float32(analysisInst.AnalysisTemperature),
	}
}
