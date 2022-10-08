package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/beacon_indexer/beacon_models"
)

const InsertForwardFetchCheckpointUsingCache = "InsertForwardFetchCheckpointUsingCache"

var UpdateUsingCacheTimeout = time.Minute * 5

func UpdateAllValidatorBalancesFromCache() {
	log.Info().Msg("UpdateAllValidatorBalancesFromCache")
	for {
		timeBegin := time.Now()
		err := fetchAllValidatorsToUpdateFromCache(context.Background(), UpdateUsingCacheTimeout)
		log.Err(err)
		log.Info().Interface("UpdateAllValidatorBalancesFromCache took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(3 * time.Second)
	}
}

func fetchAllValidatorsToUpdateFromCache(ctx context.Context, contextTimeout time.Duration) error {
	log.Info().Msg("fetchAllValidatorsToUpdateFromCache")
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	err := Fetcher.InsertForwardFetchCheckpointUsingCache(ctxTimeout)
	log.Info().Err(err).Msg("UpdateAllValidators: fetchAllValidatorsToUpdate")
	return err
}

func (f *BeaconFetcher) InsertForwardFetchCheckpointUsingCache(ctx context.Context) error {
	log.Info().Msg(InsertForwardFetchCheckpointUsingCache)

	chkPoint := beacon_models.ValidatorsEpochCheckpoint{}
	err := chkPoint.GetsOrderedNextEpochCheckpointWithBalancesRemainingAfterEpoch(ctx, checkpointEpoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalancesAfterCheckpoint: GetsOrderedNextEpochCheckpointWithBalancesRemainingAfterEpoch")
		return err
	}
	log.Info().Err(err).Msgf(InsertForwardFetchCheckpointUsingCache+" : updating at epoch %d", chkPoint.Epoch)
	vbe, err := Fetcher.FetchForwardCheckpointValidatorBalances(ctx, int64(chkPoint.Epoch))
	if err != nil {
		log.Err(err).Msg(InsertForwardFetchCheckpointUsingCache)
		return err
	}
	err = vbe.InsertValidatorBalances(ctx)
	if err != nil {
		log.Err(err).Msg(InsertForwardFetchCheckpointUsingCache)
		return err
	}
	err = beacon_models.UpdateEpochCheckpointBalancesRecordedAtEpoch(ctx, chkPoint.Epoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorsToUpdateFromCache: " + InsertForwardFetchCheckpointUsingCache)
		return err
	}
	return nil
}
