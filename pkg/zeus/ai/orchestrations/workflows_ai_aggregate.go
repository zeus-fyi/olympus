package ai_platform_service_orchestrations

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"go.temporal.io/sdk/workflow"
)

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
	ao := getDefaultRetryPolicy()
	i := runCycle

	for _, aggInst := range wfExecParams.WorkflowTasks {
		log.Info().Interface("runCycle", runCycle).Msg("aggregation: runCycle")
		if isInvalidAggInst(aggInst, md, wfExecParams) {
			continue
		}
		aggCycle := wfExecParams.CycleCountTaskRelative.AggNormalizedCycleCounts[*aggInst.AggTaskID]
		// checks for run cycle validity
		if i%aggCycle != 0 {
			continue
		}
		logger.Info("aggregation: taskID", *aggInst.AggTaskID)
		cp.Tc = getAggTaskContext(aggInst)
		cp.Window = artemis_orchestrations.CalculateTimeWindowFromCycles(wfExecParams.WorkflowExecTimekeepingParams.RunWindow.UnixStartTime, i-aggCycle, i, wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize)
		aggRetrievalCtx := workflow.WithActivityOptions(ctx, ao)
		err := workflow.ExecuteActivity(aggRetrievalCtx, z.AiAggregateAnalysisRetrievalTask, cp, getAnalysisDeps(aggInst, wfExecParams)).Get(aggRetrievalCtx, &cp)
		if err != nil {
			logger.Error("failed to run aggregate retrieval", "Error", err)
			return err
		}
		md.AggregateAnalysis[*aggInst.AggTaskID][aggInst.AnalysisTaskID] = false
		log.Info().Msg("aggregation: running token overflow reduction")
		var chunkIterator int
		chunkedTaskCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(chunkedTaskCtx, z.TokenOverflowReduction, cp, nil).Get(chunkedTaskCtx, &cp)
		if err != nil {
			logger.Error("failed to run agg token overflow reduction task", "Error", err)
			return err
		}
		chunkIterator = cp.Tc.ChunkIterator
		log.Info().Interface("chunkIterator", chunkIterator).Msg("agg: chunkIterator")
		for chunkOffset := 0; chunkOffset < chunkIterator; chunkOffset++ {
			log.Info().Interface("chunkOffset", chunkOffset).Msg("agg: chunkOffset")
			cp.Wsr.ChunkOffset = chunkOffset
			switch aws.StringValue(aggInst.AggResponseFormat) {
			case jsonFormat:
				log.Info().Interface("jsonFormat", jsonFormat).Msg("agg: jsonFormat")
				childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{WorkflowID: oj.OrchestrationName + "-agg-json-task-" + strconv.Itoa(i), WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize, RetryPolicy: ao.RetryPolicy}
				childAggWfCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
				cp.Wsr.ChildWfID = childAnalysisWorkflowOptions.WorkflowID
				err = workflow.ExecuteChildWorkflow(childAggWfCtx, z.JsonOutputTaskWorkflow, cp).Get(childAggWfCtx, &cp)
				if err != nil {
					logger.Error("failed to execute json agg workflow", "Error", err)
					return err
				}
			case csvFormat:
				aofoa := ao
				aofoa.HeartbeatTimeout = 5 * time.Minute
				aggCompCtx := workflow.WithActivityOptions(ctx, aofoa)
				err = workflow.ExecuteActivity(aggCompCtx, z.CsvIterator, cp).Get(aggCompCtx, nil)
				if err != nil {
					logger.Error("failed to save agg csv response", "Error", err)
					return err
				}
			case ingestData:
				aofoa := ao
				aofoa.HeartbeatTimeout = 10 * time.Minute
				aggCompCtx := workflow.WithActivityOptions(ctx, aofoa)
				err = workflow.ExecuteActivity(aggCompCtx, z.RecordCompletionResponse, ou, getDummyChatCompResp()).Get(aggCompCtx, &cp.Tc.ResponseID)
				if err != nil {
					logger.Error("failed to save agg read only response", "Error", err)
					return err
				}
				wr := getWr(cp, chunkOffset)
				cp.Tc.ResponseFormat = csvFormat
				recordAggCtx := workflow.WithActivityOptions(ctx, aofoa)
				err = workflow.ExecuteActivity(recordAggCtx, z.SaveCsvTaskOutput, cp, wr).Get(recordAggCtx, &cp.Tc.WorkflowResultID)
				if err != nil {
					logger.Error("failed to save agg analysis", "Error", err)
					return err
				}
			default:
				aofoa := ao
				aofoa.HeartbeatTimeout = 10 * time.Minute
				log.Info().Interface("textFormat", text).Msg("agg: textFormat")
				var aiAggResp *ChatCompletionQueryResponse
				aggCtx := workflow.WithActivityOptions(ctx, aofoa)
				err = workflow.ExecuteActivity(aggCtx, z.AiAggregateTask, aggInst, cp).Get(aggCtx, &aiAggResp)
				if err != nil {
					logger.Error("failed to run aggregation", "Error", err)
					return err
				}
				aggCompCtx := workflow.WithActivityOptions(ctx, aofoa)
				err = workflow.ExecuteActivity(aggCompCtx, z.RecordCompletionResponse, ou, aiAggResp).Get(aggCompCtx, &cp.Tc.ResponseID)
				if err != nil {
					logger.Error("failed to save agg response", "Error", err)
					return err
				}
				wr := getWr(cp, chunkOffset)
				ia := InputDataAnalysisToAgg{
					ChatCompletionQueryResponse: aiAggResp,
				}
				var tmp string
				for _, cv := range aiAggResp.Response.Choices {
					tmp += cv.Message.Content + "\n"
				}
				ia.TextInput = &tmp
				recordAggCtx := workflow.WithActivityOptions(ctx, aofoa)
				err = workflow.ExecuteActivity(recordAggCtx, z.SaveTaskOutput, wr, cp, ia).Get(recordAggCtx, &cp.Tc.WorkflowResultID)
				if err != nil {
					logger.Error("failed to save aggregation resp", "Error", err)
					return err
				}
			}
			for ind, evalFn := range aggInst.AggEvalFns {
				evalAggCycle := wfExecParams.CycleCountTaskRelative.AggEvalNormalizedCycleCounts[*aggInst.AggTaskID][evalFn.EvalID]
				if i%evalAggCycle == 0 {
					childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
						WorkflowID:         oj.OrchestrationName + "-agg-eval-" + strconv.Itoa(i) + "-chunk-" + strconv.Itoa(chunkOffset) + "-ind-" + strconv.Itoa(ind),
						RetryPolicy:        ao.RetryPolicy,
						WorkflowRunTimeout: ao.ScheduleToCloseTimeout,
					}
					cp.Tc.EvalID = evalFn.EvalID
					log.Info().Interface("evalFn.EvalID", evalFn.EvalID).Msg("agg: eval")
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
		log.Info().Interface("len(pr.PromptReductionText.OutPromptChunks)", len(pr.PromptReductionText.OutPromptChunks)).Msg("getChunkIteratorLen: PromptReductionText")
		return len(pr.PromptReductionText.OutPromptChunks)
	}
	return chunkIterator
}

func getAggTaskContext(aggInst artemis_orchestrations.WorkflowTemplateData) TaskContext {
	return TaskContext{
		TaskName:              aws.StringValue(aggInst.AggTaskName),
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
}
