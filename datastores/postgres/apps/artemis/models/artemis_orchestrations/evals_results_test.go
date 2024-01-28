package artemis_orchestrations

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func (s *OrchestrationsTestSuite) TestUpsertEvalMetricsResults() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	eocr := EvalMetaDataResult{
		EvalOpCtxStr: "zzzz",
	}
	b, err := json.Marshal(eocr)
	s.Require().Nil(err)
	evCtx := EvalContext{
		EvalID:             1705978298687209000,
		EvalIterationCount: 0,
		AIWorkflowAnalysisResult: AIWorkflowAnalysisResult{
			WorkflowResultID: 1706407535130182000,
			OrchestrationID:  1706407459064898000,
			//ResponseID:            1672188679693780000,
			SourceTaskID:          1705949866538066000,
			IterationCount:        0,
			ChunkOffset:           0,
			RunningCycleNumber:    1,
			SearchWindowUnixStart: 1706400471,
			SearchWindowUnixEnd:   1706400772,
			SkipAnalysis:          false,
			Metadata:              nil,
			CompletionChoices:     nil,
		},
	}
	emrw := &EvalMetricsResults{
		EvalContext: evCtx,
		EvalMetricsResults: []*EvalMetric{
			{
				EvalMetricID: aws.Int(1706135988888304000),
				EvalMetricResult: &EvalMetricResult{
					EvalMetricResultID:    aws.Int(1706407544885429000),
					EvalResultOutcomeBool: aws.Bool(true),
					EvalMetadata:          b,
				},
				EvalOperator:            "gt",
				EvalState:               "info",
				EvalExpectedResultState: "pass",
				EvalMetricComparisonValues: &EvalMetricComparisonValues{
					EvalComparisonNumber: aws.Float64(322),
				},
			},
			{
				EvalMetricID: aws.Int(1706151746655067000),
				EvalMetricResult: &EvalMetricResult{
					EvalMetricResultID:    aws.Int(1706151746655067006),
					EvalResultOutcomeBool: aws.Bool(true),
					EvalMetadata:          b,
				},
				EvalOperator:            "gt",
				EvalState:               "filter",
				EvalExpectedResultState: "fail",
				EvalMetricComparisonValues: &EvalMetricComparisonValues{
					EvalComparisonNumber: aws.Float64(556),
				},
			},
			{
				EvalMetricID: aws.Int(1706135988888304000),
				EvalMetricResult: &EvalMetricResult{
					EvalMetricResultID:    aws.Int(1706151746655067009),
					EvalResultOutcomeBool: aws.Bool(true),
					EvalMetadata:          b,
				},
				EvalOperator:            "gt",
				EvalState:               "info",
				EvalExpectedResultState: "pass",
				EvalMetricComparisonValues: &EvalMetricComparisonValues{
					EvalComparisonNumber: aws.Float64(322),
				},
			},
		},
	}
	err = UpsertEvalMetricsResults(ctx, emrw)
	s.Require().Nil(err)

	fp := filepaths.Path{
		DirIn: "/Users/alex/go/Olympus/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations",
		FnIn:  "emrs_debug.json",
	}

	b = fp.ReadFileInPath()

	var emr EvalMetricsResults
	err = json.Unmarshal(b, &emr)
	s.Require().Nil(err)

	m := make(map[int]int)

	for _, v := range emr.EvalMetricsResults {
		fmt.Println(*v.EvalMetricResult.EvalMetricResultID)
		m[*v.EvalMetricResult.EvalMetricResultID]++
	}
	for k, v := range m {
		fmt.Println(k, v)
	}

	err = UpsertEvalMetricsResults(ctx, &emr)
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
