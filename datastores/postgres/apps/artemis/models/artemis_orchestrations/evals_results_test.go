package artemis_orchestrations

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (s *OrchestrationsTestSuite) TestUpsertEvalMetricsResults() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	evCtx := EvalContext{
		EvalID:             1705978298687209000,
		EvalIterationCount: 0,
		AIWorkflowAnalysisResult: AIWorkflowAnalysisResult{
			WorkflowResultID: 1705978298687209000,
			OrchestrationsID: 1705978283783333000,
			ResponseID:       1672188679693780000,
			SourceTaskID:     1705949866538066000,
			IterationCount:   0,
			ChunkOffset:      0,
			SkipAnalysis:     false,
		},
	}
	emrw := &EvalMetricsResults{
		EvalContext: evCtx,
		EvalMetricsResults: []*EvalMetric{
			{
				EvalMetricID: aws.Int(1706151746655067000),
				EvalMetricResult: &EvalMetricResult{
					EvalResultOutcomeBool: aws.Bool(true),
				},
				EvalOperator:            "gt",
				EvalState:               "info",
				EvalExpectedResultState: "pass",
				EvalMetricComparisonValues: &EvalMetricComparisonValues{
					EvalComparisonNumber: aws.Float64(3),
				},
			},
		},
	}
	err := UpsertEvalMetricsResults(ctx, emrw)
	s.Require().Nil(err)

}

func (s *OrchestrationsTestSuite) TestInsertOrUpdateAiWorkflowEvalResultResponse() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	zz := AIWorkflowEvalResultResponse{
		EvalResultID:     1706230901688275000,
		WorkflowResultID: 1705978298687209000,
		ResponseID:       1672188679693780000,
	}
	_, err := InsertOrUpdateAiWorkflowEvalResultResponse(ctx, zz)
	s.Require().Nil(err)
}
