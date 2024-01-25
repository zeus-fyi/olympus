package ai_platform_service_orchestrations

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

/*
type EvalFn struct {
	EvalID         *int                          `json:"evalID,omitempty"`
	EvalType       string                        `json:"evalType"`
	EvalModel      *string                       `json:"evalModel,omitempty"`
	EvalFormat     string                        `json:"evalFormat"`
	EvalCycleCount int                           `json:"evalCycleCount,omitempty"`
	TriggerActions []TriggerAction               `json:"triggerFunctions,omitempty"`
	Schemas        []*JsonSchemaDefinition       `json:"schemas,omitempty"`
	SchemasMap     map[int]*JsonSchemaDefinition `json:"schemaMap"`
}
type JsonSchemaDefinition struct {
	SchemaID    int               `db:"schema_id" json:"schemaID"`
	IsObjArray  bool              `db:"is_obj_array" json:"isObjArray"`
	Fields      []JsonSchemaField `db:"-" json:"fields"`
}

type JsonSchemaField struct {
	FieldID           int         `db:"field_id" json:"fieldID"`
	FieldName         string      `db:"field_name" json:"fieldName"`
	FieldDescription  string      `db:"field_description" json:"fieldDescription"`
	DataType          string      `db:"data_type" json:"dataType"`
	IntegerValue      *int        `db:"-" json:"intValue,omitempty"`
	StringValue       *string     `db:"-" json:"stringValue,omitempty"`
	NumberValue       *float64    `db:"-" json:"numberValue,omitempty"`
	BooleanValue      *bool       `db:"-" json:"booleanValue,omitempty"`
	IntegerValueSlice []int       `db:"-" json:"intValueSlice,omitempty"`
	StringValueSlice  []string    `db:"-" json:"stringValueSlice,omitempty"`
	NumberValueSlice  []float64   `db:"-" json:"numberValueSlice,omitempty"`
	BooleanValueSlice []bool      `db:"-" json:"booleanValueSlice,omitempty"`
	IsValidated       bool        `db:"-" json:"isValidated,omitempty"`
	EvalMetric        *EvalMetric `db:"-" json:"evalMetricResult,omitempty"`
}

type EvalMetricResult struct {
	EvalMetricResultID    *int            `json:"evalMetricResultID"`
	EvalResultOutcomeBool *bool           `json:"evalResultOutcome,omitempty"` // true if eval passed, false if eval failed
	EvalMetadata          json.RawMessage `json:"evalMetadata,omitempty"`
}
type EvalMetric struct {
	EvalMetricID            *int              `json:"evalMetricID"`
	EvalMetricResult        *EvalMetricResult `json:"evalMetricResult"`
	EvalOperator            string            `json:"evalOperator"`
	EvalState               string            `json:"evalState"`
	EvalExpectedResultState string            `json:"evalExpectedResultState"` // true if eval passed, false if eval failed
	EvalComparisonBoolean   *bool             `json:"evalComparisonBoolean,omitempty"`
	EvalComparisonNumber    *float64          `json:"evalComparisonNumber,omitempty"`
	EvalComparisonString    *string           `json:"evalComparisonString,omitempty"`
	EvalComparisonInteger   *int              `json:"evalComparisonInteger,omitempty"`
	EvalMetadata            json.RawMessage   `json:"evalMetadata,omitempty"`
}
*/

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
					}
				case "msg_id":
					for _, em := range field.EvalMetrics {
						t.Require().NotNil(em.EvalMetricResult)
						t.Require().NotNil(em.EvalMetricResult.EvalResultOutcomeBool)
						t.Require().True(*em.EvalMetricResult.EvalResultOutcomeBool)
						count += 1000
					}
					// Add additional cases for other field types as needed
				}
			}
		}
	}
	t.Require().Equal(1111, count, "Expected 11 fields to be validated")

}
