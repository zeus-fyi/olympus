package ai_platform_service_orchestrations

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (t *ZeusWorkerTestSuite) TestJsonToEvalMetric() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)

	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID
	act := NewZeusAiPlatformActivities()

	//for _, evalFn := range evalFnMetrics {
	//	for _, metric := range evalFn.SchemasMap {
	//		err = TransformJSONToEvalScoredMetrics(metric)
	//		if err != nil {
	//			t.Require().Nil(err)
	//		}
	//	}
	//}

	td, err := act.SelectTaskDefinition(ctx, ou, 1701313525731432000)
	t.Require().Nil(err)
	t.Require().NotEmpty(td)
	taskSchemaMap := make(map[string]*artemis_orchestrations.JsonSchemaDefinition)

	for ti, _ := range td {
		for si, _ := range td[ti].Schemas {
			for fi, _ := range td[ti].Schemas[si].Fields {
				if td[ti].Schemas[si].Fields[fi].FieldName == "msg_score" {
					td[ti].Schemas[si].Fields[fi].IntegerValue = aws.Int(2)
					td[ti].Schemas[si].Fields[fi].IsValidated = true
				}
				if td[ti].Schemas[si].Fields[fi].FieldName == "msg_id" {
					td[ti].Schemas[si].Fields[fi].IntegerValue = aws.Int(1111111111111)
					td[ti].Schemas[si].Fields[fi].IsValidated = true
				}
			}
			taskSchemaMap[td[ti].Schemas[si].SchemaStrID] = td[ti].Schemas[si]
		}
	}

	evalFnMetrics, err := act.EvalLookup(ctx, ou, 1706151746524505000)
	t.Require().Nil(err)
	t.Require().NotEmpty(evalFnMetrics)
	for ei, evalFn := range evalFnMetrics {
		if evalFn.EvalID == nil {
			continue
		}
		copyMatchingFieldValues(taskSchemaMap, evalFnMetrics[ei].SchemasMap)
	}

	count := 0
	// Now, verify that the values have been transferred correctly
	for _, evalFn := range evalFnMetrics {
		if evalFn.EvalID == nil {
			continue
		}
		for _, schema := range evalFn.SchemasMap {
			for _, field := range schema.Fields {
				// Check if the field values have been copied correctly based on field names
				switch field.FieldName {
				case "msg_score":
					t.Require().NotNil(field.IntegerValue, "IntegerValue should not be nil for 'score'")
					t.Require().Equal(2, *field.IntegerValue, "IntegerValue should be 2 for 'score'")
					t.Require().True(field.IsValidated, "IsValidated should be true for 'score'")
					count += 1
				case "msg_id":
					t.Require().NotNil(field.IntegerValue, "IntegerValue should not be nil for 'msg_id'")
					t.Require().Equal(1111111111111, *field.IntegerValue, "IntegerValue should be not 0 for 'msg_id'")
					t.Require().True(field.IsValidated, "IsValidated should be true for 'msg_id'")
					count += 10
					// Add additional cases for other field types as needed
				}
			}
			err = TransformJSONToEvalScoredMetrics(schema)
			t.Require().Nil(err)

			for _, field := range schema.Fields {
				switch field.FieldName {
				case "msg_score":
					for _, em := range field.EvalMetrics {
						t.Require().NotNil(em.EvalMetricResult)
						t.Require().NotNil(em.EvalMetricResult.EvalResultOutcomeBool)
						t.Require().False(*em.EvalMetricResult.EvalResultOutcomeBool)
						count += 100

						eorc := artemis_orchestrations.EvalMetaDataResult{}
						err = json.Unmarshal(em.EvalMetricResult.EvalMetadata, &eorc)
						t.Require().Nil(err)
						t.Require().NotNil(eorc)
						t.Require().NotNil(eorc.EvalComparisonInteger)
					}
				case "msg_id":
					for _, em := range field.EvalMetrics {
						t.Require().NotNil(em.EvalMetricResult)
						t.Require().NotNil(em.EvalMetricResult.EvalResultOutcomeBool)
						t.Require().True(*em.EvalMetricResult.EvalResultOutcomeBool)
						count += 1000
						eorc := artemis_orchestrations.EvalMetaDataResult{}
						err = json.Unmarshal(em.EvalMetricResult.EvalMetadata, &eorc)
						t.Require().NotNil(eorc.EvalComparisonInteger)
					}
					// Add additional cases for other field types as needed
				}
			}
		}
	}
	t.Require().Equal(1111, count, "Expected 11 fields to be validated")

}
