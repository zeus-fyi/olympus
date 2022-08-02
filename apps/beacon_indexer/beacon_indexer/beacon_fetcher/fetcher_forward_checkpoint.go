package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/beacon_api"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres/beacon_indexer/beacon_models"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

var FetchAllValidatorBalancesTimeoutFromCheckpoint = time.Minute * 10
var checkpointEpoch = 134000

func FetchAllValidatorBalancesAfterCheckpoint() {
	log.Info().Msg("FetchFindAndQueryAndUpdateValidatorBalances")

	for {
		timeBegin := time.Now()
		err := fetchAllValidatorBalancesAfterCheckpoint(context.Background(), FetchAllValidatorBalancesTimeoutFromCheckpoint)
		log.Err(err)
		log.Info().Interface("FetchFindAndQueryAndUpdateValidatorBalances took this many seconds to complete: ", time.Now().Sub(timeBegin))
	}
}

func fetchAllValidatorBalancesAfterCheckpoint(ctx context.Context, contextTimeout time.Duration) error {
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
	err = chkPoint.GetAnyEpochCheckpointWithBalancesRemainingAfterEpoch(ctx, checkpointEpoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalances")
		return err
	}
	log.Info().Msgf("Fetching balances for all active validators at epoch %d", chkPoint.Epoch)

	if fetcher.Cache.DoesCheckpointExist(ctx, chkPoint.Epoch) {
		log.Info().Msgf("Fetching balances skipping api call since, checkpoint cache exists at epoch %d", chkPoint.Epoch)
		return nil
	}

	balances, err := fetcher.FetchAllValidatorBalances(ctxTimeout, int64(chkPoint.Epoch))
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

	fetcher.Cache.SetCheckpointCache(ctx, chkPoint.Epoch, 10*time.Minute)
	log.Info().Err(err).Msg("fetchAllValidatorBalances")
	return err
}

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
