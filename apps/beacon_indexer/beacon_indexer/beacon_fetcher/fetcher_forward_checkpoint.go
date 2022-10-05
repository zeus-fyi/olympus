package beacon_fetcher

import (
	"context"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/beacon_indexer/beacon_models"
	"github.com/zeus-fyi/olympus/pkg/beacon_api"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

var FetchAllValidatorBalancesTimeoutFromCheckpoint = time.Minute * 5
var FetchAnyValidatorBalancesTimeoutFromCheckpoint = time.Minute * 3
var checkpointEpoch = 134000

func FetchAllValidatorBalancesAfterCheckpoint() {
	log.Info().Msg("FetchAllValidatorBalancesAfterCheckpoint")

	for {
		timeBegin := time.Now()
		err := fetchAllValidatorBalancesAfterCheckpoint(context.Background(), FetchAllValidatorBalancesTimeoutFromCheckpoint)
		log.Err(err)
		log.Info().Interface("FetchFindAndQueryAndUpdateValidatorBalances took this many seconds to complete: ", time.Now().Sub(timeBegin))
	}
}

func FetchAnyValidatorBalancesAfterCheckpoint() {
	log.Info().Msg("FetchAnyValidatorBalancesAfterCheckpoint")

	for {
		timeBegin := time.Now()
		err := fetchAnyValidatorBalancesAfterCheckpoint(context.Background(), FetchAnyValidatorBalancesTimeoutFromCheckpoint)
		log.Err(err)
		log.Info().Interface("fetchAnyValidatorBalancesAfterCheckpoint took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(30 * time.Second)
	}
}
func fetchAnyValidatorBalancesAfterCheckpoint(ctx context.Context, contextTimeout time.Duration) error {
	log.Info().Msg("fetchAnyValidatorBalancesAfterCheckpoint")
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	chkPoint := beacon_models.ValidatorsEpochCheckpoint{}
	err := chkPoint.GetsOrderedNextEpochCheckpointWithBalancesRemainingAfterEpoch(ctx, checkpointEpoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalancesAfterCheckpoint: GetFirstEpochCheckpointWithBalancesRemaining")
		return err
	}
	min := 2
	max := 100
	findEpoch := (rand.Intn(max-min+1) + min) + chkPoint.Epoch

	err = chkPoint.GetAnyEpochCheckpointWithBalancesRemainingAfterEpoch(ctx, findEpoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAnyValidatorBalancesAfterCheckpoint: GetAnyEpochCheckpointWithBalancesRemainingAfterEpoch")
		return err
	}
	log.Info().Msgf("fetchAnyValidatorBalancesAfterCheckpoint: Fetching balances for all active validators at epoch %d", findEpoch)

	if isCached, cacheErr := Fetcher.Cache.DoesCheckpointExist(ctx, findEpoch); cacheErr != nil {
		log.Error().Err(cacheErr).Msg("fetchAnyValidatorBalancesAfterCheckpoint: DoesCheckpointExist")
	} else if isCached {
		log.Info().Msgf("fetchAnyValidatorBalancesAfterCheckpoint: skipping fetch balance api call since, checkpoint cache exists at epoch %d", findEpoch)
	}

	_, err = Fetcher.FetchAllValidatorBalances(ctxTimeout, int64(findEpoch))
	if err != nil {
		log.Info().Err(err).Msgf("fetchAnyValidatorBalancesAfterCheckpoint: FetchAllValidatorBalances at Epoch: %d", findEpoch)
		return err
	}
	return err
}

func fetchAllValidatorBalancesAfterCheckpoint(ctx context.Context, contextTimeout time.Duration) error {
	log.Info().Msg("fetchAllValidatorBalancesAfterCheckpoint")
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	chkPoint := beacon_models.ValidatorsEpochCheckpoint{}
	err := chkPoint.GetNextEpochCheckpoint(ctx)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalancesAfterCheckpoint")
	}

	err = beacon_models.UpdateEpochCheckpointBalancesRecordedAtEpoch(ctxTimeout, chkPoint.Epoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalancesAfterCheckpoint: UpdateEpochCheckpointBalancesRecordedAtEpoch")
	}
	err = chkPoint.GetsOrderedNextEpochCheckpointWithBalancesRemainingAfterEpoch(ctx, checkpointEpoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalancesAfterCheckpoint")
		return err
	}
	log.Info().Msgf("fetchAllValidatorBalancesAfterCheckpoint: Fetching balances for all active validators at epoch %d", chkPoint.Epoch)

	if isCached, cacheErr := Fetcher.Cache.DoesCheckpointExist(ctx, chkPoint.Epoch); cacheErr != nil {
		log.Error().Err(cacheErr).Msg("fetchAllValidatorBalancesAfterCheckpoint: DoesCheckpointExist")
	} else if isCached {
		log.Info().Msgf("fetchAllValidatorBalancesAfterCheckpoint: skipping fetch balance api call since, checkpoint cache exists at epoch %d", chkPoint.Epoch)
	}

	balances, err := Fetcher.FetchAllValidatorBalances(ctxTimeout, int64(chkPoint.Epoch))
	if err != nil {
		log.Info().Err(err).Msgf("fetchAllValidatorBalancesAfterCheckpoint: FetchAllValidatorBalances at Epoch: %d", chkPoint.Epoch)
		return err
	}
	err = balances.InsertValidatorBalancesForNextEpoch(ctx)
	if err != nil {
		log.Error().Err(err).Msg("fetchAllValidatorBalancesAfterCheckpoint: InsertValidatorBalancesForNextEpoch")
		return err
	}

	err = beacon_models.UpdateEpochCheckpointBalancesRecordedAtEpoch(ctx, chkPoint.Epoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalancesAfterCheckpoint: UpdateEpochCheckpointBalancesRecordedAtEpoch")
		return err
	}

	key, err := Fetcher.Cache.SetCheckpointCache(ctx, chkPoint.Epoch, 1*time.Minute)
	log.Info().Err(err).Msgf("fetchAllValidatorBalancesAfterCheckpoint: set key failed %s", key)
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

// Checkpoints
// UpdateForwardEpochCheckpoint // Routine FOUR
func UpdateForwardEpochCheckpoint() {
	log.Info().Msg("UpdateForwardEpochCheckpoint")
	for {
		timeBegin := time.Now()
		err := checkpointForwardUpdater(context.Background(), UpdateCheckpointsTimeout)
		log.Err(err)
		log.Info().Interface("UpdateForwardEpochCheckpoint took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(UpdateCheckpointsTimeout)
	}
}

func checkpointForwardUpdater(ctx context.Context, contextTimeout time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()
	chkPoint := beacon_models.ValidatorsEpochCheckpoint{}
	err := chkPoint.GetsOrderedNextEpochCheckpointWithBalancesRemainingAfterEpoch(ctx, checkpointEpoch)
	if err != nil {
		log.Info().Err(err).Msg("checkpointForwardUpdater")
		return err
	}
	log.Info().Msgf("UpdateForwardEpochCheckpoint: checkpointForwardUpdater at Epoch %d", chkPoint.Epoch)
	err = beacon_models.UpdateEpochCheckpointBalancesRecordedAtEpoch(ctxTimeout, chkPoint.Epoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalances: checkpointUpdater")
		return err
	}
	return err
}
