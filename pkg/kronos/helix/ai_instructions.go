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
