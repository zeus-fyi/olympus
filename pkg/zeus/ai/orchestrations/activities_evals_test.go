package ai_platform_service_orchestrations

import (
	"fmt"

	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (t *ZeusWorkerTestSuite) TestEvalFnJsonConstructSchemaFromDb() {
	apps.Pg.InitPG(ctx, t.Tc.LocalDbPgconn)

	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID

	act := NewZeusAiPlatformActivities()
	evalFnMetrics, err := act.EvalLookup(ctx, ou, 1703624059411640000)
	t.Require().Nil(err)
	t.Require().NotEmpty(evalFnMetrics)

	jd, err := TransformEvalMetricsToJSONSchema(evalFnMetrics[0].EvalMetrics)
	t.Require().Nil(err)
	t.Require().NotEmpty(jd)

	fdSchema := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"count": {
				Type:        jsonschema.Number,
				Description: "total number of words in sentence",
			},
			"words": {
				Type:        jsonschema.Array,
				Description: "list of words in sentence",
				Items: &jsonschema.Definition{
					Type: jsonschema.String,
				},
			},
		},
		Required: []string{"count", "words"},
	}
	t.Require().Equal(fdSchema.Type, jd.Type)
	t.Require().Equal(fdSchema.Properties["count"].Description, jd.Properties["count"].Description)
	t.Require().Equal(fdSchema.Properties["count"].Type, jd.Properties["count"].Type)
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

	jsonData := `{"count": 10, "words": ["word1", "word2"]}`

	metrics, err := TransformJSONToEvalScoredMetrics(jsonData, evalFnMetrics[0].EvalMetricMap)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Output the resulting metrics
	for _, metric := range metrics {
		fmt.Printf("Metric: %+v\n", metric)
	}
}
