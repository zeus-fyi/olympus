package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/beacon_indexer/beacon_models"
)

var UpdateCheckpointsTimeout = 60 * time.Second
var InsertCheckpointsTimeout = time.Minute * 1

// InsertNewEpochCheckpoint // Routine FIVE
func InsertNewEpochCheckpoint() {
	log.Info().Msg("InsertNewEpochCheckpoint")
	for {
		timeBegin := time.Now()
		err := newCheckpoint(context.Background(), InsertCheckpointsTimeout)
		log.Err(err)
		log.Info().Interface("InsertNewEpochCheckpoint took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(InsertCheckpointsTimeout)
	}
}

func newCheckpoint(ctx context.Context, contextTimeout time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()
	chkPoint := beacon_models.ValidatorsEpochCheckpoint{}
	err := chkPoint.GetNextEpochCheckpoint(ctx)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalances")
		return err
	}
	log.Info().Msgf("InsertNewEpochCheckpoint: newCheckpoint at Epoch %d", chkPoint.Epoch)
	finalizedEpoch := beacon_models.ValidatorsEpochCheckpoint{}

	err = finalizedEpoch.GetCurrentFinalizedEpoch(ctx)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalances: GetCurrentFinalizedEpoch")
		return err
	}

	if chkPoint.Epoch > finalizedEpoch.Epoch {
		log.Info().Msg("fetchAllValidatorBalances: GetCurrentFinalizedEpoch, checkpoint epoch is up to date with finalized")
		time.Sleep(time.Minute * 3)
		return nil
	}
	_, err = beacon_models.InsertEpochCheckpoint(ctxTimeout, chkPoint.Epoch)
	if err != nil {
		log.Error().Err(err).Msg("InsertNewEpochCheckpoint: newCheckpoint")
	}
	err = beacon_models.UpdateEpochCheckpointBalancesRecordedAtEpoch(ctxTimeout, chkPoint.Epoch)
	if err != nil {
		log.Error().Err(err).Msg("InsertNewEpochCheckpoint: newCheckpoint")
	}
	return err
}
