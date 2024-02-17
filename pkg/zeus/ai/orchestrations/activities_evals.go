package ai_platform_service_orchestrations

import (
	"context"

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

func (z *ZeusAiPlatformActivities) SaveEvalMetricResults(ctx context.Context, emr *artemis_orchestrations.EvalMetricsResults) error {
	if emr == nil || emr.EvalMetricsResults == nil {
		log.Info().Msg("SaveEvalMetricResults: emr is nil")
		return nil
	}
	err := artemis_orchestrations.UpsertEvalMetricsResults(ctx, emr)
	if err != nil {
		log.Err(err).Msg("EvalLookup: failed to get eval fn")
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

func (z *ZeusAiPlatformActivities) EvalModelScoredJsonOutput(ctx context.Context, ef *artemis_orchestrations.EvalFn, jrs []artemis_orchestrations.JsonSchemaDefinition) (*artemis_orchestrations.EvalMetricsResults, error) {
	if ef == nil || ef.SchemasMap == nil {
		log.Info().Msg("EvalModelScoredJsonOutput: at least one input is nil or empty")
		return nil, nil
	}
	emr := &artemis_orchestrations.EvalMetricsResults{EvalMetricsResults: []*artemis_orchestrations.EvalMetric{}}
	for i, _ := range jrs {
		_, ok := ef.SchemasMap[jrs[i].SchemaStrID]
		if !ok {
			continue
		}
		copyMatchingFieldValuesFromResp(&jrs[i], ef.SchemasMap)
		err := TransformJSONToEvalScoredMetrics(ef.SchemasMap[jrs[i].SchemaStrID])
		if err != nil {
			log.Err(err).Int("i", i).Interface("scoredResultsFields", jrs[i]).Msg("EvalModelScoredJsonOutput: failed to transform json to eval scored metrics")
			return nil, err
		}
		for fj, _ := range ef.SchemasMap[jrs[i].SchemaStrID].Fields {
			for emi, _ := range ef.SchemasMap[jrs[i].SchemaStrID].Fields[fj].EvalMetrics {
				if ef.SchemasMap[jrs[i].SchemaStrID].Fields[fj].EvalMetrics[emi] == nil {
					continue
				}
				evm := *ef.SchemasMap[jrs[i].SchemaStrID].Fields[fj].EvalMetrics[emi]
				jrs[i].ScoredEvalMetrics = append(jrs[i].ScoredEvalMetrics, &evm)
				emr.EvalMetricsResults = append(emr.EvalMetricsResults, &evm)
			}
		}
	}
	emr.EvaluatedJsonResponses = jrs
	return emr, nil
}
