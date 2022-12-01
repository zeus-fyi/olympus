package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/beacon_indexer/beacon_models"
	"github.com/zeus-fyi/olympus/pkg/apollo/ethereum/consensus_client_apis/beacon_api"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos/v0"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

var FetchAllValidatorBalancesTimeout = time.Minute * 10

// FetchAllValidatorBalances Routine THREE
func FetchAllValidatorBalances() {
	log.Info().Msg("FetchFindAndQueryAndUpdateValidatorBalances")

	for {
		timeBegin := time.Now()
		err := fetchAllValidatorBalances(context.Background(), FetchAllValidatorBalancesTimeout)
		log.Err(err)
		log.Info().Interface("FetchFindAndQueryAndUpdateValidatorBalances took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(FetchAllValidatorBalancesTimeout)
	}
}

func fetchAllValidatorBalances(ctx context.Context, contextTimeout time.Duration) error {
	log.Info().Msg("fetchAllValidatorBalances")
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	chkPoint := beacon_models.ValidatorsEpochCheckpoint{}
	err := chkPoint.GetNextEpochCheckpoint(ctx)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalances")
	}

	err = beacon_models.UpdateEpochCheckpointBalancesRecordedAtEpoch(ctxTimeout, chkPoint.Epoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalances: UpdateEpochCheckpointBalancesRecordedAtEpoch")
	}
	err = chkPoint.GetFirstEpochCheckpointWithBalancesRemaining(ctx)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalances")
		return err
	}
	log.Info().Msgf("Fetching balances for all active validators at epoch %d", chkPoint.Epoch)

	isCached, err := Fetcher.Cache.DoesCheckpointExist(ctx, chkPoint.Epoch)
	log.Info().Err(err).Msg("fetchAllValidatorBalances: DoesCheckpointExist")
	if isCached {
		log.Info().Msgf("Fetching balances skipping api call since, checkpoint cache exists at epoch %d", chkPoint.Epoch)
		return nil
	}

	balances, err := Fetcher.FetchAllValidatorBalances(ctxTimeout, int64(chkPoint.Epoch))
	if err != nil {
		log.Info().Err(err).Msgf("fetchAllValidatorBalances: FetchAllValidatorBalances at Epoch: %d", chkPoint.Epoch)
		return err
	}
	err = balances.InsertValidatorBalancesForNextEpoch(ctx)
	if err != nil {
		log.Error().Err(err).Msg("fetchAllValidatorBalances: InsertValidatorBalancesForNextEpoch")
		return err
	}

	err = beacon_models.UpdateEpochCheckpointBalancesRecordedAtEpoch(ctx, chkPoint.Epoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalances: UpdateEpochCheckpointBalancesRecordedAtEpoch")
		return err
	}
	key, err := Fetcher.Cache.SetCheckpointCache(ctx, chkPoint.Epoch, 1*time.Minute)
	log.Info().Err(err).Msgf("fetchAllValidatorBalances: set key failed %s", key)
	return err
}

func (f *BeaconFetcher) FetchAllValidatorBalances(ctx context.Context, epoch int64) (beacon_models.ValidatorBalancesEpoch, error) {
	log.Info().Msg("BeaconFetcher: FetchAllValidatorBalancesAtSlot")
	var valBalances beacon_models.ValidatorBalancesEpoch
	var beaconAPI beacon_api.ValidatorBalances

	vbe, err := Fetcher.Cache.GetBalanceCache(ctx, int(epoch))
	if err != nil || len(vbe.ValidatorBalances) == 0 {
		log.Err(err).Msg("balance cache not found, fetching from beacon")
	} else {
		return vbe, err
	}
	lib := v0.LibV0{}
	slotToQuery := lib.ConvertEpochToSlot(epoch)
	err = beaconAPI.FetchAllValidatorBalancesAtStateAndDecode(ctx, f.NodeEndpoint, slotToQuery)
	if err != nil {
		log.Error().Err(err).Msg("BeaconFetcher: QueryAllValidatorBalancesAtSlot")
		return valBalances, err
	}
	log.Info().Msg("BeaconFetcher: Convert API data to model format")
	valBalances.ValidatorBalances = make([]beacon_models.ValidatorBalanceEpoch, len(beaconAPI.Data))
	for i, vbFromAPI := range beaconAPI.Data {
		vbForDataEntry := beacon_models.ValidatorBalanceEpoch{
			Validator:        beacon_models.Validator{Index: string_utils.Int64StringParser(vbFromAPI.Index)},
			Epoch:            epoch,
			TotalBalanceGwei: string_utils.Int64StringParser(vbFromAPI.Balance),
		}
		valBalances.ValidatorBalances[i] = vbForDataEntry
	}

	_, err = Fetcher.Cache.SetBalanceCache(ctx, int(epoch), valBalances, time.Hour*24*1)
	if err != nil {
		log.Err(err)
	}
	return valBalances, nil
}
