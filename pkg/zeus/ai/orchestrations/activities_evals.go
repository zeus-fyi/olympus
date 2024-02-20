package ai_platform_service_orchestrations

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (z *ZeusAiPlatformActivities) SendResponseToApiForScoresInJson(ctx context.Context, cr *ChatCompletionQueryResponse) (map[string]interface{}, error) {
	return nil, nil
}

func (z *ZeusAiPlatformActivities) EvalLookup(ctx context.Context, ou org_users.OrgUser, evalID int) ([]artemis_orchestrations.EvalFn, error) {
	evalFn, err := artemis_orchestrations.SelectEvalFnsByOrgIDAndID(ctx, ou, evalID)
	if err != nil {
		log.Err(err).Msg("EvalLookup: failed to get eval fn")
		return nil, err
	}
	return evalFn, nil
}

func gws(ctx context.Context, inputID int) (WorkflowStageIO, error) {
	act := NewZeusAiPlatformActivities()
	wio, err := act.SelectWorkflowIO(ctx, inputID)
	if err != nil {
		log.Err(err).Msg("SaveEvalResponseOutput: failed to select workflow io")
		return wio, err
	}
	return wio, err
}

func sws(ctx context.Context, input *WorkflowStageIO) (*WorkflowStageIO, error) {
	if input == nil {
		log.Warn().Msg("SaveEvalResponseOutput: at least one input is nil or empty")
		return nil, nil
	}
	act := NewZeusAiPlatformActivities()
	wio, err := act.SaveWorkflowIO(ctx, input)
	if err != nil {
		log.Err(err).Msg("SaveEvalResponseOutput: failed to select workflow io")
		return wio, err
	}
	return wio, err
}

func (z *ZeusAiPlatformActivities) SaveEvalMetricResults(ctx context.Context, inputID int) error {
	wio, err := gws(ctx, inputID)
	if err != nil {
		log.Err(err).Msg("SaveEvalMetricResults: failed to select workflow io")
		return err
	}
	if wio.CreateTriggerActionsWorkflowInputs == nil || wio.CreateTriggerActionsWorkflowInputs.Emr == nil {
		log.Info().Msg("SaveEvalMetricResults: at least one input is nil or empty")
		return nil
	}
	wio.Logs = append(wio.Logs, "EvalMetricResults started")
	err = artemis_orchestrations.UpsertEvalMetricsResults(ctx, wio.CreateTriggerActionsWorkflowInputs.Emr)
	if err != nil {
		log.Err(err).Msg("EvalLookup: failed to get eval fn")
		return err
	}
	wio.Logs = append(wio.Logs, "EvalMetricResults done")
	_, err = sws(ctx, &wio)
	if err != nil {
		log.Err(err).Msg("SaveEvalMetricResults: failed to save workflow io")
		return err
	}
	return nil
}

func (z *ZeusAiPlatformActivities) SaveEvalResponseOutput(ctx context.Context, errr artemis_orchestrations.AIWorkflowEvalResultResponse) (int, error) {
	respID, err := artemis_orchestrations.InsertOrUpdateAiWorkflowEvalResultResponse(ctx, errr)
	if err != nil {
		log.Err(err).Interface("respID", respID).Interface("errr", errr).Msg("SaveEvalResponseOutput: failed")
		return -1, err
	}
	return respID, nil
}

func (z *ZeusAiPlatformActivities) EvalModelScoredJsonOutput(ctx context.Context, ef *artemis_orchestrations.EvalFn, cp *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	if ef == nil || ef.SchemasMap == nil || cp == nil {
		log.Info().Msg("EvalModelScoredJsonOutput: at least one input is nil or empty")
		return nil, nil
	}
	act := NewZeusAiPlatformActivities()
	wio, err := act.SelectWorkflowIO(ctx, cp.Wsr.InputID)
	if err != nil {
		log.Err(err).Msg("SaveEvalResponseOutput: failed to select workflow io")
		return nil, err
	}
	emr := &artemis_orchestrations.EvalMetricsResults{
		EvalContext: artemis_orchestrations.EvalContext{
			EvalID: aws.ToInt(ef.EvalID),
			AIWorkflowAnalysisResult: artemis_orchestrations.AIWorkflowAnalysisResult{
				OrchestrationID:       cp.Oj.OrchestrationID,
				SourceTaskID:          cp.Tc.TaskID,
				IterationCount:        cp.Wsr.IterationCount,
				ChunkOffset:           cp.Wsr.ChunkOffset,
				RunningCycleNumber:    cp.Wsr.RunCycle,
				SearchWindowUnixStart: cp.Window.UnixStartTime,
				SearchWindowUnixEnd:   cp.Window.UnixEndTime,
				ResponseID:            cp.Tc.ResponseID,
			},
			AIWorkflowEvalResultResponse: artemis_orchestrations.AIWorkflowEvalResultResponse{
				EvalResultsID:      cp.Tc.EvalResultID,
				EvalID:             aws.ToInt(ef.EvalID),
				WorkflowResultID:   cp.Tc.WorkflowResultID,
				ResponseID:         cp.Tc.ResponseID,
				EvalIterationCount: cp.Wsr.EvalIterationCount,
			},
		},
		EvalMetricsResults: []*artemis_orchestrations.EvalMetric{},
	}
	for i, _ := range cp.Tc.JsonResponseResults {
		_, ok := ef.SchemasMap[cp.Tc.JsonResponseResults[i].SchemaStrID]
		if !ok {
			continue
		}
		copyMatchingFieldValuesFromResp(&cp.Tc.JsonResponseResults[i], ef.SchemasMap)
		err = TransformJSONToEvalScoredMetrics(ef.SchemasMap[cp.Tc.JsonResponseResults[i].SchemaStrID])
		if err != nil {
			log.Err(err).Int("i", i).Interface("scoredResultsFields", cp.Tc.JsonResponseResults[i]).Msg("EvalModelScoredJsonOutput: failed to transform json to eval scored metrics")
			return nil, err
		}
		for fj, _ := range ef.SchemasMap[cp.Tc.JsonResponseResults[i].SchemaStrID].Fields {
			for emi, _ := range ef.SchemasMap[cp.Tc.JsonResponseResults[i].SchemaStrID].Fields[fj].EvalMetrics {
				if ef.SchemasMap[cp.Tc.JsonResponseResults[i].SchemaStrID].Fields[fj].EvalMetrics[emi] == nil {
					continue
				}
				evm := *ef.SchemasMap[cp.Tc.JsonResponseResults[i].SchemaStrID].Fields[fj].EvalMetrics[emi]
				cp.Tc.JsonResponseResults[i].ScoredEvalMetrics = append(cp.Tc.JsonResponseResults[i].ScoredEvalMetrics, &evm)
				emr.EvalMetricsResults = append(emr.EvalMetricsResults, &evm)
			}
		}
	}
	emr.EvaluatedJsonResponses = cp.Tc.JsonResponseResults
	if wio.CreateTriggerActionsWorkflowInputs == nil {
		wio.CreateTriggerActionsWorkflowInputs = &CreateTriggerActionsWorkflowInputs{}
	}
	wio.CreateTriggerActionsWorkflowInputs.Emr = emr
	err = artemis_orchestrations.UpsertEvalMetricsResults(ctx, wio.CreateTriggerActionsWorkflowInputs.Emr)
	if err != nil {
		log.Err(err).Msg("EvalLookup: failed to get eval fn")
		return nil, err
	}
	_, err = act.SaveWorkflowIO(ctx, &wio)
	if err != nil {
		log.Err(err).Msg("SaveEvalResponseOutput: failed to save workflow io")
		return nil, err
	}
	cp.Tc.JsonResponseResults = emr.EvaluatedJsonResponses
	return cp, nil
}
