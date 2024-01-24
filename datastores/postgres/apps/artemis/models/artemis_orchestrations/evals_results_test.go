package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (s *OrchestrationsTestSuite) TestUpsertEvalMetricsResults() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	evCtx := EvalContext{
		EvalID:             0,
		EvalIterationCount: 0,
		AIWorkflowAnalysisResult: AIWorkflowAnalysisResult{
			WorkflowResultID: 0,
			OrchestrationsID: 0,
			ResponseID:       0,
			SourceTaskID:     0,
			IterationCount:   0,
			ChunkOffset:      0,
			SkipAnalysis:     false,
		},
	}
	emrw := EvalMetricsResults{
		EvalContext: evCtx,
		EvalMetricsResults: []EvalMetric{
			{
				EvalMetricResult:        nil,
				EvalOperator:            "",
				EvalState:               "",
				EvalExpectedResultState: "",
				EvalMetricComparisonValues: &EvalMetricComparisonValues{
					EvalComparisonBoolean: nil,
					EvalComparisonNumber:  nil,
					EvalComparisonString:  nil,
					EvalComparisonInteger: nil,
				},
			},
		},
	}
	err := UpsertEvalMetricsResults(ctx, emrw.EvalContext, emrw.EvalMetricsResults)
	s.Require().Nil(err)
	zz := AIWorkflowEvalResultResponse{
		EvalResultID:     0,
		WorkflowResultID: 0,
		EvalID:           0,
		ResponseID:       0,
	}
	_, err = InsertOrUpdateAiWorkflowEvalResultResponse(ctx, zz)
	s.Require().Nil(err)
}

func (s *OrchestrationsTestSuite) TestInsertOrUpdateAiWorkflowEvalResultResponse() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
}
