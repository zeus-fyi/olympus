package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres/beacon_indexer/beacon_models"
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
	err := beacon_models.UpdateEpochCheckpointBalancesRecordedAtEpoch(ctxTimeout, chkPoint.Epoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalances: UpdateEpochCheckpointBalancesRecordedAtEpoch")
	}
	err = chkPoint.GetFirstEpochCheckpointWithBalancesRemaining(ctx)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalances")
		return err
	}
	log.Info().Msgf("Fetching balances for all active validators at epoch %d", chkPoint.Epoch)

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

	log.Info().Err(err).Msg("fetchAllValidatorBalances")
	return err
}
