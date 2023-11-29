package kronos_helix

import "time"

type AiInstructions struct {
	WorkflowName string                   `json:"workflowName"`
	StepSize     int                      `json:"stepSize"`
	StepSizeUnit string                   `json:"stepSizeUnit"`
	StartTime    time.Time                `json:"startTime"`
	Map          map[string]AiInstruction `json:"map"`
}

type AiInstruction struct {
	InstructionType       string `json:"instructionType"`
	Model                 string `json:"model"`
	Prompt                string `json:"prompt"`
	MaxTokens             int    `json:"maxTokens"`
	TokenOverflowStrategy string `json:"tokenOverflowStrategy"`
	CycleCount            int    `json:"cycleCount"`

	// working variables
	TokensConsumed int `json:"tokensConsumed"`
	CyclesLeft     int `json:"cyclesLeft"`
}

func (a *AiInstructions) GetInstructionByType(aiInstType string) AiInstruction {
	tmp := a.Map["analysis"]
	switch aiInstType {
	case "analysis":
		tmp = a.Map["analysis"]
	case "aggregation":
		tmp = a.Map["aggregation"]
	}
	return tmp
}

//type AiSearchParams struct {
//	SearchContentText                     string       `json:"searchContentText,omitempty"`
//	GroupFilter                           string       `json:"groupFilter,omitempty"`
//	Platforms                             string       `json:"platforms,omitempty"`
//	Usernames                             string       `json:"usernames,omitempty"`
//	WorkflowInstructions                  string       `json:"workflowInstructions,omitempty"`
//	WorkflowCycleInstructions             string       `json:"workflowCycleInstructions,omitempty"`
//	SearchInterval                        TimeInterval `json:"searchInterval,omitempty"`
//	AnalysisInterval                      TimeInterval `json:"analysisInterval,omitempty"`
//	CycleCount                            int          `json:"cycleCount,omitempty"`
//	AggregationCycleCount                 int          `json:"aggregationCycleCount,omitempty"`
//	StepSize                              int          `json:"stepSize,omitempty"`
//	StepSizeUnit                          string       `json:"stepSizeUnit,omitempty"`
//	TimeRange                             string       `json:"timeRange,omitempty"`
//	WorkflowName                          string       `json:"workflowName,omitempty"`
//	AnalysisModel                         string       `json:"analysisModel,omitempty"`
//	AnalysisModelMaxTokens                int          `json:"analysisModelMaxTokens,omitempty"`
//	AnalysisModelTokenOverflowStrategy    string       `json:"analysisModelTokenOverflowStrategy,omitempty"`
//	AggregationModel                      string       `json:"aggregationModel,omitempty"`
//	AggregationModelMaxTokens             int          `json:"aggregationModelMaxTokens,omitempty"`
//	AggregationModelTokenOverflowStrategy string       `json:"aggregationModelTokenOverflowStrategy,omitempty"`
//}
