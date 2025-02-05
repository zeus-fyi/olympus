package artemis_orchestrations

import (
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
						FieldID:          0,
						FieldName:        "",
						FieldDescription: "",
						DataType:         "number",
						FieldValue: FieldValue{
							IntegerValue:      nil,
							StringValue:       nil,
							NumberValue:       nil,
							BooleanValue:      nil,
							IntegerValueSlice: nil,
							StringValueSlice:  nil,
							NumberValueSlice:  nil,
							BooleanValueSlice: nil,
							IsValidated:       false,
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
						FieldID:          0,
						FieldName:        "",
						FieldDescription: "",
						DataType:         "array[string]",
						FieldValue:       FieldValue{},
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

func (s *OrchestrationsTestSuite) TestInsertEvalWithJsonSchema() {
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
					},
					{
						FieldID:          1705808653789520000,
						FieldName:        "lead_score_metrics",
						FieldDescription: "lead_score_metrics description",
						DataType:         "array[string]",
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
						FieldID:          0,
						FieldName:        "",
						FieldDescription: "",
						DataType:         "number",
						FieldValue: FieldValue{
							IntegerValue:      nil,
							StringValue:       nil,
							NumberValue:       nil,
							BooleanValue:      nil,
							IntegerValueSlice: nil,
							StringValueSlice:  nil,
							NumberValueSlice:  nil,
							BooleanValueSlice: nil,
							IsValidated:       false,
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
						FieldID:          0,
						FieldName:        "",
						FieldDescription: "",
						DataType:         "array[string]",
						FieldValue:       FieldValue{},
					},
				},
			},
		},
	}
	err := InsertOrUpdateEvalFnWithMetrics(ctx, ou, &evalFn)
	s.Require().Nil(err)
	s.Require().NotNil(evalFn.EvalID)
}

//
//func (s *OrchestrationsTestSuite) TestInsertEvalResults() {
//	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
//	// Example data for EvalFn and EvalMetrics
//	ou := org_users.OrgUser{}
//	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
//	ou.UserID = s.Tc.ProductionLocalTemporalUserID
//	// Example data for EvalContext
//	evalContext := EvalContext{
//		EvalID:             0,
//		EvalIterationCount: 0,
//		AIWorkflowAnalysisResult: AIWorkflowAnalysisResult{
//			WorkflowResultID:      0,
//			OrchestrationsID:      0,
//			ResponseID:            0,
//			SourceTaskID:          0,
//			IterationCount:        0,
//			ChunkOffset:           0,
//			RunningCycleNumber:    0,
//			SearchWindowUnixStart: 0,
//			SearchWindowUnixEnd:   0,
//			SkipAnalysis:          false,
//		},
//	}
//
//	// Example data for EvalMetricsResults
//	evalMetricsResults := []EvalMetric{
//		{
//			EvalMetricID: aws.Int(1702959527792954000), // Assume a valid eval metric ID
//			//EvalMetadata: json.RawMessage(`{"example": "metadata1"}`),
//		},
//		{
//			EvalMetricID: aws.Int(1703624059422669000), // Assume another valid eval metric I
//			//EvalMetadata: json.RawMessage(`{"example": "metadata2"}`),
//		},
//	}
//
//	// Call the function to insert data
//	err := UpsertEvalMetricsResults(ctx, evalContext, evalMetricsResults)
//	s.Require().NoError(err)
//}
