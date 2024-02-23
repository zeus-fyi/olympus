package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/beacon_indexer/beacon_models"
	"github.com/zeus-fyi/olympus/pkg/apollo/ethereum/client_apis/beacon_api"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos/v0"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func (f *BeaconFetcher) FetchAllValidatorBalances(ctx context.Context, epoch int) (beacon_models.ValidatorBalancesEpoch, error) {
	log.Info().Msg("BeaconFetcher: FetchAllValidatorBalancesAtSlot")
	var valBalances beacon_models.ValidatorBalancesEpoch
	var beaconAPI beacon_api.ValidatorBalances

	vbe, err := Fetcher.Cache.GetBalanceCache(ctx, epoch)
	if err != nil || len(vbe.Data) == 0 {
		log.Err(err).Msg("balance cache not found, fetching from beacon")
	}

	if len(vbe.Data) == 0 {
		lib := v0.LibV0{}
		slotToQuery := lib.ConvertEpochToSlot(int64(epoch))
		b, berr := beaconAPI.FetchAllValidatorBalancesAtStateAndDecode(ctx, f.NodeEndpoint, slotToQuery)
		if berr != nil {
			log.Err(berr).Msg("BeaconFetcher: QueryAllValidatorBalancesAtSlot")
			return valBalances, berr
		}
		b.Epoch = epoch
		_, cerr := Fetcher.Cache.SetBalanceCache(ctx, epoch, b, time.Hour*1)
		if cerr != nil {
			log.Err(cerr)
		}
		valBalances.ValidatorBalances = make([]beacon_models.ValidatorBalanceEpoch, len(beaconAPI.Data))
		for i, vbFromAPI := range beaconAPI.Data {
			vbForDataEntry := beacon_models.ValidatorBalanceEpoch{
				Validator:        beacon_models.Validator{Index: string_utils.Int64StringParser(vbFromAPI.Index)},
				Epoch:            int64(epoch),
				TotalBalanceGwei: string_utils.Int64StringParser(vbFromAPI.Balance),
			}
			valBalances.ValidatorBalances[i] = vbForDataEntry
		}
		return valBalances, nil
	}
	log.Info().Msg("BeaconFetcher: Convert API data to model format")

	valBalances.ValidatorBalances = make([]beacon_models.ValidatorBalanceEpoch, len(vbe.Data))
	for i, vbFromAPI := range vbe.Data {
		vbForDataEntry := beacon_models.ValidatorBalanceEpoch{
			Validator:        beacon_models.Validator{Index: string_utils.Int64StringParser(vbFromAPI.Index)},
			Epoch:            int64(epoch),
			TotalBalanceGwei: string_utils.Int64StringParser(vbFromAPI.Balance),
		}
		valBalances.ValidatorBalances[i] = vbForDataEntry
	}
	return valBalances, nil
}
