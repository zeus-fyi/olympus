package api_types

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/beacon-indexer/beacon_indexer/beacon_api"
	"github.com/zeus-fyi/olympus/pkg/client"
)

type ValidatorBalances struct {
	Epoch int64

	ExecutionOptimistic bool `json:"execution_optimistic"`
	Data                []struct {
		Index   string `json:"index"`
		Balance string `json:"balance"`
	} `json:"data"`
}

func (b *ValidatorBalances) FetchStateAndDecode(ctx context.Context, beaconNode, stateID string, encodedQueryURL string) error {
	log.Info().Msg("ValidatorBalances: FetchStateAndDecode")
	r := beacon_api.GetValidatorsBalancesByStateFilter(ctx, beaconNode, stateID, encodedQueryURL)

	log.Info().Interface("GetValidatorsBalancesByStateFilter: Status Code Response: ", r.Status)
	if r.Err != nil {
		log.Error().Err(r.Err).Msg("FetchStateAndDecode: FetchStateAndDecode")
	}

	return b.DecodeValidatorsBalancesBeacon(r)
}

func (b *ValidatorBalances) FetchAllValidatorBalancesAtStateAndDecode(ctx context.Context, beaconNode, stateID string) error {
	log.Info().Msg("ValidatorBalances: FetchAllValidatorBalancesAtStateAndDecode")
	r := beacon_api.GetAllValidatorBalancesByState(ctx, beaconNode, stateID)

	if r.StatusCode != http.StatusOK {
		log.Info().Interface("FetchAllValidatorBalancesAtStateAndDecode: Status Code Response: ", r.Status)
		return errors.New("request had an unexpected non-200 status code response")
	}

	if r.Err != nil {
		log.Error().Err(r.Err).Msg("FetchStateAndDecode: FetchAllValidatorBalancesAtStateAndDecode")
	}

	return b.DecodeValidatorsBalancesBeacon(r)
}
func (b *ValidatorBalances) DecodeValidatorsBalancesBeacon(r client.Reply) error {
	log.Info().Msg("ValidatorBalances: DecodeValidatorsBalancesBeacon")
	err := json.Unmarshal(r.BodyBytes, &b)

	if err != nil {
		log.Info().Str("DecodeValidatorsBalancesBeacon STATUS CODE ", r.Status)
		log.Info().Interface("ValidatorBalances: DecodeValidatorsBalancesBeacon ", &b)
		log.Error().Err(err).Msg("ValidatorBalances: DecodeValidatorsBalancesBeacon")
	}
	return err
}
