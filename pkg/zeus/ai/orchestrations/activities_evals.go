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

func (z *ZeusAiPlatformActivities) SaveEvalResponseOutput(ctx context.Context, errr artemis_orchestrations.AIWorkflowEvalResultResponse) error {
	respID, err := artemis_orchestrations.InsertOrUpdateAiWorkflowEvalResultResponse(ctx, errr)
	if err != nil {
		log.Err(err).Interface("respID", respID).Interface("errr", errr).Msg("SaveTaskOutput: failed")
		return err
	}
	return nil
}

//
//func (z *ZeusAiPlatformActivities) EvalFnModelScoredJsonOutput(ctx context.Context, schemas [][]*artemis_orchestrations.JsonSchemaDefinition) (*artemis_orchestrations.EvalMetricsResults, error) {
//	if schemas == nil {
//		log.Info().Msg("EvalModelScoredJsonOutput: at least one input is nil or empty")
//		return nil, nil
//	}
//	emr := &artemis_orchestrations.EvalMetricsResults{
//		EvalMetricsResults: []*artemis_orchestrations.EvalMetric{},
//	}
//	for i, _ := range schemas {
//		for j, _ := range schemas[i] {
//			err := TransformJSONToEvalScoredMetrics(schemas[i][j])
//			if err != nil {
//				log.Err(err).Int("i", i).Interface("scoredResultsFields", schemas[i][j]).Msg("EvalModelScoredJsonOutput: failed to transform json to eval scored metrics")
//				return nil, err
//			}
//
//			for fj, _ := range schemas[i][j].Fields {
//				for emi, _ := range schemas[i][j].Fields[fj].EvalMetrics {
//					emr.EvalMetricsResults = append(emr.EvalMetricsResults, schemas[i][j].Fields[fj].EvalMetrics[emi])
//				}
//			}
//		}
//	}
//	return emr, nil
//}

func (z *ZeusAiPlatformActivities) EvalModelScoredJsonOutput(ctx context.Context, ef *artemis_orchestrations.EvalFn, jrs [][]*artemis_orchestrations.JsonSchemaDefinition) (*artemis_orchestrations.EvalMetricsResults, error) {
	if ef == nil || ef.SchemasMap == nil {
		log.Info().Msg("EvalModelScoredJsonOutput: at least one input is nil or empty")
		return nil, nil
	}

	emr := &artemis_orchestrations.EvalMetricsResults{
		EvalMetricsResults: []*artemis_orchestrations.EvalMetric{},
	}

	for i, _ := range jrs {
		for j, _ := range jrs[i] {
			tmp := jrs[i][j]
			_, ok := ef.SchemasMap[tmp.SchemaStrID]
			if !ok {
				continue
			}
			copyMatchingFieldValuesFromResp(tmp, ef.SchemasMap)
			err := TransformJSONToEvalScoredMetrics(ef.SchemasMap[tmp.SchemaStrID])
			if err != nil {
				log.Err(err).Int("i", i).Interface("scoredResultsFields", jrs[i][j]).Msg("EvalModelScoredJsonOutput: failed to transform json to eval scored metrics")
				return nil, err
			}
			for fj, _ := range ef.SchemasMap[tmp.SchemaStrID].Fields {
				for emi, _ := range ef.SchemasMap[tmp.SchemaStrID].Fields[fj].EvalMetrics {
					emr.EvalMetricsResults = append(emr.EvalMetricsResults, ef.SchemasMap[tmp.SchemaStrID].Fields[fj].EvalMetrics[emi])
				}
			}
		}
	}
	return emr, nil
}
