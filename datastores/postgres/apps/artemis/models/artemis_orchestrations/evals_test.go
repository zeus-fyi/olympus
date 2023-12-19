package artemis_orchestrations

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

func (s *OrchestrationsTestSuite) TestInsertEval() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	// Example data for EvalFn and EvalMetrics
	evalModel := "ExampleModel"
	evalMetricName := "Metric1"
	evalComparisonBoolean := true
	evalComparisonNumber := 42

	evalFn := EvalFn{
		EvalID:        nil, // nil implies a new entry
		OrgID:         s.Tc.ProductionLocalTemporalOrgID,
		UserID:        s.Tc.ProductionLocalTemporalUserID,
		EvalName:      "TestEvalFunction",
		EvalType:      "TypeA",
		EvalGroupName: "GroupX",
		EvalModel:     &evalModel,
		EvalFormat:    "Format1",
		EvalMetrics: []EvalMetric{
			{
				EvalModelPrompt:       "Prompt1",
				EvalMetricName:        evalMetricName,
				EvalMetricResult:      "Result1",
				EvalComparisonBoolean: &evalComparisonBoolean,
				EvalComparisonNumber:  &evalComparisonNumber,
				EvalComparisonString:  nil, // Testing with a nil value
				EvalMetricDataType:    "DataType1",
				EvalOperator:          "Operator1",
				EvalState:             "State1",
			},
			// Add more EvalMetric entries here to cover different scenarios
		},
	}

	err := InsertOrUpdateEvalFnWithMetrics(ctx, &evalFn)
	s.Require().Nil(err)
	s.Require().NotNil(evalFn.EvalID)
}

func (s *OrchestrationsTestSuite) TestSelectEvals() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	// Example data for EvalFn and EvalMetrics

	evs, err := SelectEvalFnsByOrgID(ctx, s.Tc.ProductionLocalTemporalOrgID)
	s.Require().Nil(err)
	s.Require().NotNil(evs)
}