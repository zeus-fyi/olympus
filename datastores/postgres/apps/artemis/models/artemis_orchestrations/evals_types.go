package artemis_orchestrations

import "encoding/json"

type EvalFn struct {
	EvalStrID      *string                          `json:"evalStrID,omitempty"`
	EvalID         *int                             `json:"evalID,omitempty"`
	OrgID          int                              `json:"orgID,omitempty"`
	UserID         int                              `json:"userID,omitempty"`
	EvalName       string                           `json:"evalName"`
	EvalType       string                           `json:"evalType"`
	EvalGroupName  string                           `json:"evalGroupName"`
	EvalModel      *string                          `json:"evalModel,omitempty"`
	EvalFormat     string                           `json:"evalFormat"`
	EvalCycleCount int                              `json:"evalCycleCount,omitempty"`
	TriggerActions []TriggerAction                  `json:"triggerFunctions,omitempty"`
	Schemas        []*JsonSchemaDefinition          `json:"schemas,omitempty"`
	SchemasMap     map[string]*JsonSchemaDefinition `json:"schemaMap"`
}

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

type EvalMetricResult struct {
	EvalMetricResultStrID     *string `json:"evalMetricsResultStrID"`
	EvalMetricResultID        *int    `json:"evalMetricsResultID"`
	EvalResultOutcomeBool     *bool   `json:"evalResultOutcomeBool,omitempty"`     // true if eval passed, false if eval failed, nil if ignored
	EvalResultOutcomeStateStr *string `json:"evalResultOutcomeStateStr,omitempty"` // pass, fail, ignore
	RunningCycleNumber        *int    `json:"runningCycleNumber,omitempty"`
	EvalIterationCount        *int    `json:"evalIterationCount,omitempty"`

	SearchWindowUnixStart *int            `json:"searchWindowUnixStart,omitempty"`
	SearchWindowUnixEnd   *int            `json:"searchWindowUnixEnd,omitempty"`
	EvalMetadata          json.RawMessage `json:"evalMetadata,omitempty"`
}

/*
		EvalTriggerState     string `db:"eval_trigger_state" json:"evalTriggerState"` // eg. info, filter, etc
		EvalResultsTriggerOn string `db:"eval_results_trigger_on" json:"evalResultsTriggerOn"` // all-pass, any-fail, etc

		this contains the result of all-pass, any-fail, etc per each element json item vs all elements
	    so we want to only use the passing elements regardless of the trigger on, since that is already accounted for
		on a per element level
*/

type EvalMetaDataResult struct {
	EvalOpCtxStr string `json:"evalOpCtxStr"`
	Operator     string `json:"operator"`
	*EvalMetricComparisonValues
	*FieldValue
}

type EvalMetric struct {
	EvalMetricStrID            *string                     `json:"evalMetricStrID,omitempty"`
	EvalMetricID               *int                        `json:"evalMetricID"`
	EvalField                  *JsonSchemaField            `json:"evalField,omitempty"`
	EvalName                   *string                     `json:"evalName,omitempty"`
	EvalMetricResult           *EvalMetricResult           `json:"evalMetricResult,omitempty"`
	EvalOperator               string                      `json:"evalOperator"`
	EvalState                  string                      `json:"evalState"`
	EvalExpectedResultState    string                      `json:"evalExpectedResultState"`
	EvalMetricComparisonValues *EvalMetricComparisonValues `json:"evalMetricComparisonValues,omitempty"`
}

type EvalMetricComparisonValues struct {
	EvalComparisonBoolean *bool    `json:"evalComparisonBoolean,omitempty"`
	EvalComparisonNumber  *float64 `json:"evalComparisonNumber,omitempty"`
	EvalComparisonString  *string  `json:"evalComparisonString,omitempty"`
	EvalComparisonInteger *int     `json:"evalComparisonInteger,omitempty"`
}
