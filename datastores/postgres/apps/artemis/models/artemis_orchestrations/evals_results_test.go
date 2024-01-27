package artemis_orchestrations

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (s *OrchestrationsTestSuite) TestUpsertEvalMetricsResults() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	eocr := EvalMetaDataResult{
		EvalOpCtxStr: "4 gt 3",
	}
	b, err := json.Marshal(eocr)
	s.Require().Nil(err)
	evCtx := EvalContext{
		EvalID:             1705978298687209000,
		EvalIterationCount: 0,
		AIWorkflowAnalysisResult: AIWorkflowAnalysisResult{
			WorkflowResultID: 1705978298687209000,
			OrchestrationID:  1705978283783333000,
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
					EvalMetricResultID:    aws.Int(1706151746655067000),
					EvalResultOutcomeBool: aws.Bool(true),
					EvalMetadata:          b,
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
	err = UpsertEvalMetricsResults(ctx, emrw)
	s.Require().Nil(err)

}

func (s *OrchestrationsTestSuite) TestInsertOrUpdateAiWorkflowEvalResultResponse() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	zz := AIWorkflowEvalResultResponse{
		EvalID:             1705978298687209000,
		WorkflowResultID:   1705978298687209000,
		ResponseID:         1672188679693780000,
		EvalIterationCount: 0,
	}
	_, err := InsertOrUpdateAiWorkflowEvalResultResponse(ctx, zz)
	s.Require().Nil(err)
}
