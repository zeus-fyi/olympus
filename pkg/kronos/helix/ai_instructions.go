package kronos_helix

import (
	"time"
)

type AiInstructions struct {
	WorkflowName string                          `json:"workflowName"`
	StepSize     int                             `json:"stepSize"`
	StepSizeUnit string                          `json:"stepSizeUnit"`
	StartTime    time.Time                       `json:"startTime"`
	Iteration    int                             `json:"iterations"`
	Map          map[string][]AiModelInstruction `json:"map"`
}

func Iterate(is *AiInstructions) []AiModelInstruction {
	is.Iteration++
	var tmp []AiModelInstruction
	ordering := []string{"analysis", "aggregation"}
	for _, o := range ordering {
		ai, ok := is.Map[o]
		if ok {
			for _, a := range ai {
				if is.Iteration%a.CycleCount == 0 {
					tmp = append(tmp, a)
				}
			}
		}
	}
	return tmp
}

type AiModelInstruction struct {
	Model                 string `json:"model"`
	Prompt                string `json:"prompt"`
	MaxTokens             int    `json:"maxTokens"`
	TokenOverflowStrategy string `json:"tokenOverflowStrategy"`
	CycleCount            int    `json:"cycleCount"`

	// working variables
	TokensConsumed int `json:"tokensConsumed"`
}
