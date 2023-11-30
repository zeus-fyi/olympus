package kronos_helix

import "fmt"

// You can change any params for this, it is a template of the other test meant for creating alerts
func (t *KronosWorkerTestSuite) TestAiOrdering() {

	ais := AiInstructions{
		WorkflowName: "test",
		StepSize:     1,
		StepSizeUnit: "minute",
		Iteration:    0,
		Map: map[string][]AiModelInstruction{
			"analysis": {
				{
					Prompt:     "prompt-analysis",
					Model:      "model",
					CycleCount: 1,
				},
			},
			"aggregation": {
				{
					Prompt:     "prompt-aggregation",
					Model:      "model",
					CycleCount: 3,
				},
			},
		},
	}

	for i := 0; i < 12; i++ {
		m := Iterate(&ais)
		fmt.Println(m)
	}
	fmt.Println(ais.Iteration)
}

/*
type AiInstructions struct {
	WorkflowName string                       `json:"workflowName"`
	StepSize     int                          `json:"stepSize"`
	StepSizeUnit string                       `json:"stepSizeUnit"`
	StartTime    time.Time                    `json:"startTime"`
	Iterations   int                          `json:"iterations"`
	Map          map[int][]AiModelInstruction `json:"map"`
	Slice        []AiModelInstruction         `json:"slice"`
}

func AggregateInstructions(is []AiModelInstruction) map[int][]AiModelInstruction {
	m := make(map[int][]AiModelInstruction)
	for i, ai := range is {
		if ai.CycleCount%i == 0 {
			m[i] = append(m[i], ai)
		}
		fmt.Println(ai.CycleCount)
	}
	return m
}

*/
