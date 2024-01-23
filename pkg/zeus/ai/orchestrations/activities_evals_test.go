package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (t *ZeusWorkerTestSuite) TestEvalFnJsonConstructSchemaFromDb() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)

	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID

	act := NewZeusAiPlatformActivities()
	evalFns, err := act.EvalLookup(ctx, ou, 1705952457940027000)
	t.Require().Nil(err)
	t.Require().NotEmpty(evalFns)

	for _, evalFnWithMetrics := range evalFns {
		fmt.Println(evalFnWithMetrics.EvalType)
		resp, rerr := EvalModelScoredJsonOutput(ctx, &evalFnWithMetrics)
		t.Require().Nil(rerr)
		t.Require().NotNil(resp)
	}
}

func EvalModelScoredJsonOutput(ctx context.Context, evalFn *artemis_orchestrations.EvalFn) (*artemis_orchestrations.EvalMetricsResults, error) {
	if evalFn == nil || evalFn.Schemas == nil {
		log.Info().Msg("EvalModelScoredJsonOutput: at least one input is nil or empty")
		return nil, nil
	}
	emr := &artemis_orchestrations.EvalMetricsResults{
		EvalMetricsResults: []artemis_orchestrations.EvalMetricsResult{},
	}
	for i, schema := range evalFn.Schemas {
		scoredResults, err := TransformJSONToEvalScoredMetrics(schema)
		if err != nil {
			log.Err(err).Int("i", i).Interface("scoredResults", scoredResults).Msg("EvalModelScoredJsonOutput: failed to transform json to eval scored metrics")
			return nil, err
		}
		emr.EvalMetricsResults = append(emr.EvalMetricsResults, scoredResults.EvalMetricsResults...)
	}
	return emr, nil
}

func (t *ZeusWorkerTestSuite) TestJsonToEvalMetric() {
	apps.Pg.InitPG(ctx, t.Tc.LocalDbPgconn)

	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID
	act := NewZeusAiPlatformActivities()
	evalFnMetrics, err := act.EvalLookup(ctx, ou, 1703624059411640000)
	t.Require().Nil(err)
	t.Require().NotEmpty(evalFnMetrics)
	// Output the resulting metrics
	//for _, evalFn := range evalFnMetrics {
	//
	//	for _, metric := range evalFn.EvalMetrics {
	//		fmt.Printf("Metric: %+v\n", metric)
	//
	//		//if metric.EvalMetricDataType == "number" {
	//		//	metric.EvalOperator = "=="
	//		//}
	//
	//		if metric.EvalMetricDataType == "array[string]" {
	//			metric.EvalComparisonString = aws.String("word1")
	//			metric.EvalOperator = "contains"
	//		}
	//	}
	//}
	//jsonData := `{"count": 10, "words": ["word1", "word2"]}`
	//metrics, err := TransformJSONToEvalScoredMetrics(jsonData, evalFnMetrics[0].EvalMetricMap)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//t.Require().NotEmpty(metrics)
	//// Output the resulting metrics
	//for _, metric := range metrics {
	//	fmt.Printf("Metric: %+v\n", metric)
	//}
}
