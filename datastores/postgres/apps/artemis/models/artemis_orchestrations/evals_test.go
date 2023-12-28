package artemis_orchestrations

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestInsertEval() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	// Example data for EvalFn and EvalMetrics
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
		EvalMetrics: []EvalMetric{
			{
				// Assuming we are creating a metric to track the 'count' property
				EvalMetricName:        "count",
				EvalMetricDataType:    "number",                            // As 'count' is a number in fdSchema
				EvalModelPrompt:       "total number of words in sentence", // Description from fdSchema
				EvalComparisonBoolean: nil,
				EvalComparisonNumber:  nil,
				EvalComparisonString:  nil,
				EvalOperator:          "", // Specific to your application logic
				EvalState:             "", // Specific to your application logic
			},
			{
				// Assuming we are creating a metric to track the 'words' property
				EvalMetricName:        "words",
				EvalMetricResult:      "",                          // This would be set based on actual evaluation results
				EvalMetricDataType:    "array[string]",             // As 'words' is an array of strings in fdSchema
				EvalModelPrompt:       "list of words in sentence", // Description from fdSchema
				EvalComparisonBoolean: nil,
				EvalComparisonNumber:  nil,
				EvalComparisonString:  nil,
				EvalOperator:          "", // Specific to your application logic
				EvalState:             "", // Specific to your application logic
			},
		},
	}
	err := InsertOrUpdateEvalFnWithMetrics(ctx, &evalFn)
	s.Require().Nil(err)
	s.Require().NotNil(evalFn.EvalID)
}

func (s *OrchestrationsTestSuite) TestSelectEvals() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	// Example data for EvalFn and EvalMetrics
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	evs, err := SelectEvalFnsByOrgIDAndID(ctx, ou, 1703624059411640000)
	s.Require().Nil(err)
	s.Require().NotNil(evs)
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
