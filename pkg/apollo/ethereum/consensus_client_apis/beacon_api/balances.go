package beacon_api

import (
	"context"

	"github.com/rs/zerolog/log"
)

type ValidatorBalances struct {
	Epoch               int  `json:"epoch"`
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Data                []struct {
		Index   string `json:"index"`
		Balance string `json:"balance"`
	} `json:"data"`
}

func (b *ValidatorBalances) FetchStateAndDecode(ctx context.Context, beaconNode, stateID string, encodedQueryURL string) (ValidatorBalances, error) {
	log.Info().Msg("ValidatorBalances: FetchStateAndDecode")
	bal, err := GetValidatorsBalancesByStateFilter(ctx, beaconNode, stateID, encodedQueryURL)
	return bal, err
}

func (b *ValidatorBalances) FetchAllValidatorBalancesAtStateAndDecode(ctx context.Context, beaconNode, stateID string) (ValidatorBalances, error) {
	log.Info().Msg("ValidatorBalances: FetchAllValidatorBalancesAtStateAndDecode")
	bal, err := GetAllValidatorBalancesByState(ctx, beaconNode, stateID)
	return bal, err
}
