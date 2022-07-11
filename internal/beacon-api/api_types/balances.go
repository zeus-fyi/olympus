package api_types

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	beacon_api "github.com/zeus-fyi/olympus/internal/beacon-api"
	"github.com/zeus-fyi/olympus/pkg/client"
)

type ValidatorBalances struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Data                []struct {
		Index   string `json:"index"`
		Balance string `json:"balance"`
	} `json:"data"`
}

func (b *ValidatorBalances) FetchStateAndDecode(ctx context.Context, beaconNode, stateID string, encodedQueryURL string) error {
	log.Info().Msg("ValidatorBalances: FetchStateAndDecode")
	r := beacon_api.GetValidatorsBalancesByStateFilter(ctx, beaconNode, stateID, encodedQueryURL)

	if r.Err != nil {
		log.Error().Err(r.Err).Msg("FetchStateAndDecode: FetchStateAndDecode")
	}

	return b.DecodeValidatorsBalancesBeacon(r)
}

func (b *ValidatorBalances) DecodeValidatorsBalancesBeacon(r client.Reply) error {
	log.Info().Msg("ValidatorBalances: DecodeValidatorsBalancesBeacon")
	err := json.Unmarshal(r.BodyBytes, &b)

	if err != nil {
		log.Info().Interface("ValidatorBalances: ", b)
		log.Error().Err(err).Msg("ValidatorBalances: DecodeValidatorsBalancesBeacon")
	}
	return err
}
