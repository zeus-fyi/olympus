package beacon_api

import (
	"context"

	"github.com/rs/zerolog/log"
)

type ValidatorsStateBeacon struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Data                []struct {
		Index     string               `json:"index"`
		Balance   string               `json:"balance"`
		Status    string               `json:"status"`
		Validator ValidatorStateBeacon `json:"validator,omitempty"`
	} `json:"data"`
}

type ValidatorStateBeacon struct {
	Pubkey                     string `json:"pubkey"`
	WithdrawalCredentials      string `json:"withdrawal_credentials"`
	EffectiveBalance           string `json:"effective_balance"`
	Slashed                    bool   `json:"slashed"`
	ActivationEligibilityEpoch string `json:"activation_eligibility_epoch"`
	ActivationEpoch            string `json:"activation_epoch"`
	ExitEpoch                  string `json:"exit_epoch"`
	WithdrawableEpoch          string `json:"withdrawable_epoch"`
}

func (b *ValidatorsStateBeacon) FetchFinalizedStateAndDecode(ctx context.Context, beaconNode string) (ValidatorsStateBeacon, error) {
	log.Info().Msg("ValidatorsStateBeacon: FetchFinalizedStateAndDecode")

	r, err := GetValidatorsFinalized(ctx, beaconNode)

	if err != nil {
		log.Error().Err(err).Msg("ValidatorsStateBeacon: FetchStateAndDecode")
	}

	return r, err
}

func (b *ValidatorsStateBeacon) FetchStateAndDecode(ctx context.Context, beaconNode, stateID, encodedQueryURL, status string) (ValidatorsStateBeacon, error) {
	log.Info().Msg("ValidatorsStateBeacon: FetchStateAndDecode")

	r, err := GetValidatorsByStateFilter(ctx, beaconNode, stateID, encodedQueryURL, status)

	if err != nil {
		log.Error().Err(err).Msg("ValidatorsStateBeacon: FetchStateAndDecode")
	}

	return r, err
}

func (b *ValidatorsStateBeacon) FetchAllStateAndDecode(ctx context.Context, beaconNode, stateID string, status string) (ValidatorsStateBeacon, error) {
	log.Info().Msg("ValidatorsStateBeacon: FetchAllStateAndDecode")

	r, err := GetValidatorsByState(ctx, beaconNode, stateID, status)

	if err != nil {
		log.Error().Err(err).Msg("ValidatorsStateBeacon: FetchAllStateAndDecode")
	}
	return r, err
}
