package api_types

type ValidatorBalances struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Data                []struct {
		Index   string `json:"index"`
		Balance string `json:"balance"`
	} `json:"data"`
}
