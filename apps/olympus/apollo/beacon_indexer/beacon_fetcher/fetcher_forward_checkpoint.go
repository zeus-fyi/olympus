package beacon_fetcher

import (
	"context"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/beacon_indexer/beacon_models"
)

var FetchAnyValidatorBalancesTimeoutFromCheckpoint = time.Minute * 1
var checkpointEpoch = 163999

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
	err := chkPoint.GetsOrderedNextEpochCheckpointWithBalancesRemainingAfterEpoch(ctxTimeout, checkpointEpoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalancesAfterCheckpoint: GetFirstEpochCheckpointWithBalancesRemaining")
		return err
	}
	minv := 2
	maxv := 20
	findEpoch := (rand.Intn(maxv-minv+1) + minv) + chkPoint.Epoch

	err = chkPoint.GetAnyEpochCheckpointWithBalancesRemainingAfterEpoch(ctxTimeout, findEpoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAnyValidatorBalancesAfterCheckpoint: GetAnyEpochCheckpointWithBalancesRemainingAfterEpoch")
		return err
	}
	log.Info().Msgf("fetchAnyValidatorBalancesAfterCheckpoint: Fetching balances for all active validators at epoch %d", findEpoch)

	if isCached, cacheErr := Fetcher.Cache.DoesCheckpointExist(ctxTimeout, findEpoch); cacheErr != nil {
		log.Err(cacheErr).Msg("fetchAnyValidatorBalancesAfterCheckpoint: DoesCheckpointExist")
	} else if isCached {
		log.Info().Msgf("fetchAnyValidatorBalancesAfterCheckpoint: skipping fetch balance api call since, checkpoint cache exists at epoch %d", findEpoch)
		return nil
	}

	// just fetch them into the cache for now... todo out of order test
	_, err = Fetcher.FetchForwardCheckpointValidatorBalances(ctx, findEpoch)
	if err != nil {
		log.Ctx(ctxTimeout).Err(err).Msg("fetchAnyValidatorBalancesAfterCheckpoint")
		return err
	}

	return err
}

func (f *BeaconFetcher) FetchForwardCheckpointValidatorBalances(ctx context.Context, epoch int) (beacon_models.ValidatorBalancesEpoch, error) {
	log.Info().Msg("BeaconFetcher: FetchForwardCheckpointValidatorBalances")
	var valBalances beacon_models.ValidatorBalancesEpoch
	// current
	vbEpochCurrent := make(map[int64]beacon_models.ValidatorBalanceEpoch)
	currentBalances, err := f.FetchAllValidatorBalances(ctx, epoch)
	if err != nil {
		log.Err(err).Msg("BeaconFetcher: FetchForwardCheckpointValidatorBalances")
		return valBalances, err
	}
	for _, vbFromAPI := range currentBalances.ValidatorBalances {
		vbEpochCurrent[vbFromAPI.Index] = vbFromAPI
	}
	// previous
	vbEpochPrevious := make(map[int64]beacon_models.ValidatorBalanceEpoch)
	prevBalances, err := f.FetchAllValidatorBalances(ctx, epoch-1)
	if err != nil {
		log.Err(err).Msg("BeaconFetcher: FetchForwardCheckpointValidatorBalances")
		return valBalances, err
	}
	for _, vbFromAPI := range prevBalances.ValidatorBalances {
		vbEpochPrevious[vbFromAPI.Index] = vbFromAPI
	}

	log.Info().Msg("BeaconFetcher: Convert API data to model format")
	valBalances.ValidatorBalances = make([]beacon_models.ValidatorBalanceEpoch, len(currentBalances.ValidatorBalances))

	i := 0
	for k, v := range vbEpochCurrent {
		currentEpochYield := 0
		if prev, ok := vbEpochPrevious[k]; ok {
			prevBalanceGwei := prev.TotalBalanceGwei
			currentEpochYield = int(v.TotalBalanceGwei - prevBalanceGwei)
		}
		vbForDataEntry := beacon_models.ValidatorBalanceEpoch{
			Validator:             beacon_models.Validator{Index: k},
			Epoch:                 int64(epoch),
			TotalBalanceGwei:      v.TotalBalanceGwei,
			CurrentEpochYieldGwei: int64(currentEpochYield),
		}
		valBalances.ValidatorBalances[i] = vbForDataEntry
		i++
	}
	return valBalances, nil
}

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
