package artemis_orchestrations

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

/*
type EvalContext struct {
	EvalID                       int                          `json:"evalID,omitempty"`
	EvalIterationCount           int                          `json:"evalIterationCount"`
	AIWorkflowAnalysisResult     AIWorkflowAnalysisResult     `json:"aiWorkflowAnalysisResult"`
	AIWorkflowEvalResultResponse AIWorkflowEvalResultResponse `json:"aiWorkflowEvalResultResponse"`
	JsonResponseResults          []JsonSchemaDefinition       `json:"jsonResponseResults,omitempty"`
	EvaluatedJsonResponses       []JsonSchemaDefinition       `json:"evaluatedJsonResponses,omitempty"`
}

type EvalMetricsResults struct {
	EvalContext        EvalContext   `json:"evalContext"`
	EvalMetricsResults []*EvalMetric `json:"evalMetricsResults"`

	EvaluatedJsonResponses []JsonSchemaDefinition `json:"evaluatedJsonResponses,omitempty"`
}
*/

func SelectEvalMetricResults(ctx context.Context, ou org_users.OrgUser) (*EvalMetricsResults, error) {
	return nil, nil
}
