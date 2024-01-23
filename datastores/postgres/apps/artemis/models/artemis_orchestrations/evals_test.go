package artemis_orchestrations

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestEvalMetricsDeleteEvalMetricsAndTriggers() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	// Example data for EvalFn and EvalMetrics
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	evalModel := "gpt-4"

	evalFn := EvalFn{
		EvalID:        aws.Int(1704860505860438000),
		OrgID:         s.Tc.ProductionLocalTemporalOrgID,
		UserID:        s.Tc.ProductionLocalTemporalUserID,
		EvalName:      "TestCountWords",
		EvalType:      "json",
		EvalGroupName: "GroupX",
		EvalModel:     &evalModel, // Assume evalModel is defined elsewhere
		EvalFormat:    "model",
		Schemas: []*JsonSchemaDefinition{
			{
				SchemaID:    0,
				SchemaName:  "",
				SchemaGroup: "",
				IsObjArray:  false,
				Fields: []JsonSchemaField{
					{
						FieldID:           0,
						FieldName:         "",
						FieldDescription:  "",
						DataType:          "number",
						IntegerValue:      nil,
						StringValue:       nil,
						NumberValue:       nil,
						BooleanValue:      nil,
						IntegerValueSlice: nil,
						StringValueSlice:  nil,
						NumberValueSlice:  nil,
						BooleanValueSlice: nil,
						EvalMetric: &EvalMetric{
							EvalMetricID:          aws.Int(1704860505915566000),
							EvalComparisonBoolean: nil,
							EvalComparisonNumber:  nil,
							EvalComparisonString:  nil,
							EvalOperator:          "", // Specific to your application logic
							EvalState:             "", // Specific to your application logic
						},
					},
				},
			},
			{
				SchemaID:    0,
				SchemaName:  "",
				SchemaGroup: "",
				IsObjArray:  false,
				Fields: []JsonSchemaField{
					{
						FieldID:           0,
						FieldName:         "",
						FieldDescription:  "",
						DataType:          "array[string]",
						IntegerValue:      nil,
						StringValue:       nil,
						NumberValue:       nil,
						BooleanValue:      nil,
						IntegerValueSlice: nil,
						StringValueSlice:  nil,
						NumberValueSlice:  nil,
						BooleanValueSlice: nil,
						EvalMetric: &EvalMetric{
							EvalMetricID:          aws.Int(1704860505860438001),
							EvalMetricResult:      "", // This would be set based on actual evaluation results
							EvalComparisonBoolean: nil,
							EvalComparisonNumber:  nil,
							EvalComparisonString:  nil,
							EvalOperator:          "", // Specific to your application logic
							EvalState:             "", // Specific to your application logic
						},
					},
				},
			},
		},
	}

	tx, err := apps.Pg.Begin(ctx)
	s.Require().Nil(err)
	defer tx.Rollback(ctx)
	evs, err := DeleteEvalMetricsAndTriggers(ctx, ou, tx, &evalFn)
	s.Require().Nil(err)
	s.Require().NotNil(evs)
	err = tx.Commit(ctx)
	s.Require().Nil(err)
}

func (s *OrchestrationsTestSuite) TestInsertEvalWithJsonSchem() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	// Example data for EvalFn and
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	evalModel := "gpt-4"
	evalFn := EvalFn{
		EvalID:        aws.Int(1705809986245485000), // nil implies a new entry
		OrgID:         s.Tc.ProductionLocalTemporalOrgID,
		UserID:        s.Tc.ProductionLocalTemporalUserID,
		EvalName:      "test-eval-schemas",
		EvalType:      "model",
		EvalGroupName: "test-eval-schemas",
		EvalModel:     &evalModel, // Assume evalModel is defined elsewhere
		EvalFormat:    "json",
		Schemas: []*JsonSchemaDefinition{
			{
				SchemaID:    1705808649277728000, // Assuming an existing schema ID for testing
				SchemaName:  "web3-sales-lead-scoring-2",
				SchemaGroup: "default",
				IsObjArray:  true,
				Fields: []JsonSchemaField{
					{
						FieldID:          1705808653782632000,
						FieldName:        "aggregate_lead_score",
						FieldDescription: "aggregate_lead_score description",
						DataType:         "number",
						EvalMetric: &EvalMetric{
							EvalMetricID:         aws.Int(1705814469816703000),
							EvalMetricResult:     "fail",
							EvalOperator:         "gt",
							EvalComparisonNumber: aws.Float64(420),
							EvalState:            "filter",
						},
					},
					{
						FieldID:          1705808653789520000,
						FieldName:        "lead_score_metrics",
						FieldDescription: "lead_score_metrics description",
						DataType:         "array[string]",
						EvalMetric: &EvalMetric{
							EvalMetricID:         aws.Int(1705812088549786000),
							EvalMetricResult:     "ignore",
							EvalOperator:         "equals",
							EvalState:            "info",
							EvalComparisonString: aws.String("test"),
						},
					},
				},
			},
			{
				SchemaID:    1705808768049976064, // Assuming an existing schema ID for testing
				SchemaName:  "messages",
				SchemaGroup: "default",
				IsObjArray:  false,
				Fields: []JsonSchemaField{
					{
						FieldID:          1705808793643712000,
						FieldName:        "msg_ids",
						FieldDescription: "system message ids",
						DataType:         "array[number]",
						EvalMetric: &EvalMetric{
							EvalMetricID:         aws.Int(1705812088550290000),
							EvalMetricResult:     "pass",
							EvalOperator:         "lt",
							EvalComparisonNumber: aws.Float64(12),
							EvalState:            "error",
						},
					},
				},
			},
		},
	}
	err := InsertOrUpdateEvalFnWithMetrics(ctx, ou, &evalFn)
	s.Require().Nil(err)
	s.Require().NotNil(evalFn.EvalID)
}

func (s *OrchestrationsTestSuite) TestInsertEval() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	// Example data for EvalFn and
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	evalModel := "gpt-4"
	evalFn := EvalFn{
		EvalID:        nil, // nil implies a new entry
		OrgID:         s.Tc.ProductionLocalTemporalOrgID,
		UserID:        s.Tc.ProductionLocalTemporalUserID,
		EvalName:      "TestCountWords",
		EvalType:      "json",
		EvalGroupName: "GroupX",
		EvalModel:     &evalModel, // Assume evalModel is defined elsewhere
		EvalFormat:    "model",
		Schemas: []*JsonSchemaDefinition{
			{
				SchemaID:    0,
				SchemaName:  "",
				SchemaGroup: "",
				IsObjArray:  false,
				Fields: []JsonSchemaField{
					{
						FieldID:           0,
						FieldName:         "",
						FieldDescription:  "",
						DataType:          "number",
						IntegerValue:      nil,
						StringValue:       nil,
						NumberValue:       nil,
						BooleanValue:      nil,
						IntegerValueSlice: nil,
						StringValueSlice:  nil,
						NumberValueSlice:  nil,
						BooleanValueSlice: nil,
						EvalMetric: &EvalMetric{
							EvalComparisonBoolean: nil,
							EvalComparisonNumber:  nil,
							EvalComparisonString:  nil,
							EvalOperator:          "", // Specific to your application logic
							EvalState:             "", // Specific to your application logic
						},
					},
				},
			},
			{
				SchemaID:    0,
				SchemaName:  "",
				SchemaGroup: "",
				IsObjArray:  false,
				Fields: []JsonSchemaField{
					{
						FieldID:           0,
						FieldName:         "",
						FieldDescription:  "",
						DataType:          "array[string]",
						IntegerValue:      nil,
						StringValue:       nil,
						NumberValue:       nil,
						BooleanValue:      nil,
						IntegerValueSlice: nil,
						StringValueSlice:  nil,
						NumberValueSlice:  nil,
						BooleanValueSlice: nil,
						EvalMetric: &EvalMetric{
							EvalMetricResult:      "", // This would be set based on actual evaluation results
							EvalComparisonBoolean: nil,
							EvalComparisonNumber:  nil,
							EvalComparisonString:  nil,
							EvalOperator:          "", // Specific to your application logic
							EvalState:             "", // Specific to your application logic
						},
					},
				},
			},
		},
	}
	err := InsertOrUpdateEvalFnWithMetrics(ctx, ou, &evalFn)
	s.Require().Nil(err)
	s.Require().NotNil(evalFn.EvalID)
}

func (s *OrchestrationsTestSuite) TestInsertEvalResults() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	// Example data for EvalFn and EvalMetrics
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	// Example data for EvalContext
	evalContext := EvalContext{
		OrchestrationID:       1701135755183363072, // Assume a valid orchestration ID
		SourceTaskID:          1701924123626391040, // Assume a valid source task ID
		RunningCycleNumber:    1,
		SearchWindowUnixStart: 1609459200, // Example Unix timestamp (e.g., Jan 1, 2021)
		SearchWindowUnixEnd:   1612137600, // Example Unix timestamp (e.g., Feb 1, 2021)
	}

	// Example data for EvalMetricsResults
	evalMetricsResults := []EvalMetricsResult{
		{
			EvalMetricID:          1702959527792954000, // Assume a valid eval metric ID
			RunningCycleNumber:    evalContext.RunningCycleNumber,
			SearchWindowUnixStart: evalContext.SearchWindowUnixStart,
			SearchWindowUnixEnd:   evalContext.SearchWindowUnixEnd,
			EvalResultOutcome:     true,
			EvalMetadata:          json.RawMessage(`{"example": "metadata1"}`),
		},
		{
			EvalMetricID:          1703624059422669000, // Assume another valid eval metric ID
			RunningCycleNumber:    evalContext.RunningCycleNumber,
			SearchWindowUnixStart: evalContext.SearchWindowUnixStart,
			SearchWindowUnixEnd:   evalContext.SearchWindowUnixEnd,
			EvalResultOutcome:     false,
			EvalMetadata:          json.RawMessage(`{"example": "metadata2"}`),
		},
	}

	// Call the function to insert data
	err := UpsertEvalMetricsResults(ctx, evalContext, evalMetricsResults)
	s.Require().NoError(err)
}
