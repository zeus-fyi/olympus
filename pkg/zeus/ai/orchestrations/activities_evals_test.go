package ai_platform_service_orchestrations

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (t *ZeusWorkerTestSuite) TestEvalFnJsonConstructSchemaFromDb() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)

	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID

	act := NewZeusAiPlatformActivities()
	evalFns, err := act.EvalLookup(ctx, ou, 1705982938987589000)
	t.Require().Nil(err)
	t.Require().NotEmpty(evalFns)

	for _, evalFnWithMetrics := range evalFns {
		fmt.Println(evalFnWithMetrics.EvalType)
		rerr := act.EvalModelScoredJsonOutput(ctx, &evalFnWithMetrics)
		t.Require().Nil(rerr)

		for _, fi := range evalFnWithMetrics.Schemas {
			for _, f := range fi.Fields {
				for _, metric := range f.EvalMetrics {
					t.Require().NotNil(metric.EvalMetricResult)
					t.Require().NotNil(metric.EvalMetricResult.EvalResultOutcomeBool)
					fmt.Println(metric.EvalMetricID, f.FieldName, *metric.EvalMetricResult.EvalResultOutcomeBool)
				}
			}
		}
	}
}

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
