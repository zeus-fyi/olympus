package api_types

type ValidatorBalances struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Data                []struct {
		Index   string `json:"index"`
		Balance string `json:"balance"`
	} `json:"data"`
}

type ValidatorStateBeacon struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Data                []struct {
		Index     string `json:"index"`
		Balance   string `json:"balance"`
		Status    string `json:"status"`
		Validator struct {
			Pubkey                     string `json:"pubkey"`
			WithdrawalCredentials      string `json:"withdrawal_credentials"`
			EffectiveBalance           string `json:"effective_balance"`
			Slashed                    bool   `json:"slashed"`
			ActivationEligibilityEpoch string `json:"activation_eligibility_epoch"`
			ActivationEpoch            string `json:"activation_epoch"`
			ExitEpoch                  string `json:"exit_epoch"`
			WithdrawableEpoch          string `json:"withdrawable_epoch"`
		} `json:"validator"`
	} `json:"data"`
}
