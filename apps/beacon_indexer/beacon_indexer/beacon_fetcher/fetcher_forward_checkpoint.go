package beacon_fetcher

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/beacon_api"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres/beacon_indexer/beacon_models"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func (f *BeaconFetcher) FetchForwardCheckpointValidatorBalances(ctx context.Context, epoch int64) (beacon_models.ValidatorBalancesEpoch, error) {
	log.Info().Msg("BeaconFetcher: FetchForwardCheckpointValidatorBalances")
	var valBalances beacon_models.ValidatorBalancesEpoch
	var beaconAPI beacon_api.ValidatorBalances

	// previous
	slotToQuery := misc.ConvertEpochToSlot(epoch - 1)
	err := beaconAPI.FetchAllValidatorBalancesAtStateAndDecode(ctx, f.NodeEndpoint, slotToQuery)
	if err != nil {
		log.Error().Err(err).Msg("BeaconFetcher: QueryAllValidatorBalancesAtSlot")
		return valBalances, err
	}
	log.Info().Msg("BeaconFetcher: Convert API data to model format")
	vbEpochPrevious := make(map[int64]beacon_models.ValidatorBalanceEpoch, len(beaconAPI.Data))
	for _, vbFromAPI := range beaconAPI.Data {
		validatorIndex := string_utils.Int64StringParser(vbFromAPI.Index)
		vbForDataEntry := beacon_models.ValidatorBalanceEpoch{
			Validator:        beacon_models.Validator{Index: validatorIndex},
			Epoch:            epoch,
			TotalBalanceGwei: string_utils.Int64StringParser(vbFromAPI.Balance),
		}
		vbEpochPrevious[validatorIndex] = vbForDataEntry
	}

	// checkpoint
	slotCheckpointToQuery := misc.ConvertEpochToSlot(epoch)
	err = beaconAPI.FetchAllValidatorBalancesAtStateAndDecode(ctx, f.NodeEndpoint, slotCheckpointToQuery)
	if err != nil {
		log.Error().Err(err).Msg("BeaconFetcher: QueryAllValidatorBalancesAtSlot")
		return valBalances, err
	}
	log.Info().Msg("BeaconFetcher: Convert API data to model format")
	valBalances.ValidatorBalances = make([]beacon_models.ValidatorBalanceEpoch, len(beaconAPI.Data))

	for i, vbFromAPI := range beaconAPI.Data {
		validatorIndex := string_utils.Int64StringParser(vbFromAPI.Index)
		totalBalance := string_utils.Int64StringParser(vbFromAPI.Balance)
		currentEpochYield := 0

		if v, ok := vbEpochPrevious[validatorIndex]; ok {
			currentEpochYield = int(totalBalance - v.TotalBalanceGwei)
		}
		vbForDataEntry := beacon_models.ValidatorBalanceEpoch{
			Validator:             beacon_models.Validator{Index: validatorIndex},
			Epoch:                 epoch,
			TotalBalanceGwei:      string_utils.Int64StringParser(vbFromAPI.Balance),
			CurrentEpochYieldGwei: int64(currentEpochYield),
		}
		valBalances.ValidatorBalances[i] = vbForDataEntry
	}

	return valBalances, nil
}
