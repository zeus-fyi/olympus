package api_types

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/beacon-indexer/beacon_indexer/beacon_api"
	"github.com/zeus-fyi/olympus/pkg/client"
)

type ValidatorsStateBeacon struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Data                []struct {
		Index     string               `json:"index"`
		Balance   string               `json:"balance"`
		Status    string               `json:"status"`
		Validator ValidatorStateBeacon `json:"validator"`
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

func (b *ValidatorsStateBeacon) FetchStateAndDecode(ctx context.Context, beaconNode, stateID string, encodedQueryURL string) error {
	r := beacon_api.GetValidatorsByStateFilter(ctx, beaconNode, stateID, encodedQueryURL)

	if r.Err != nil {
		log.Error().Err(r.Err).Msg("ValidatorsStateBeacon: FetchStateAndDecode")
	}

	return b.DecodeValidatorStateBeacon(r)
}

func (b *ValidatorsStateBeacon) FetchAllStateAndDecode(ctx context.Context, beaconNode, stateID string) error {
	r := beacon_api.GetValidatorsByState(ctx, beaconNode, stateID)

	if r.Err != nil {
		log.Error().Err(r.Err).Msg("ValidatorsStateBeacon: FetchAllStateAndDecode")
	}

	return b.DecodeValidatorStateBeacon(r)
}
func (b *ValidatorsStateBeacon) DecodeValidatorStateBeacon(r client.Reply) error {
	err := json.Unmarshal(r.BodyBytes, &b)

	if err != nil {
		log.Error().Err(err).Msg("ValidatorsStateBeacon: DecodeValidatorStateBeacon")
	}
	return err
}
