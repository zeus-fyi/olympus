package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres/beacon_indexer/beacon_models"
)

// Checkpoints

var UpdateCheckpointsTimeout = 60 * time.Second
var InsertCheckpointsTimeout = time.Minute * 2

// UpdateEpochCheckpoint // Routine FOUR
func UpdateEpochCheckpoint() {
	log.Info().Msg("UpdateEpochCheckpoint")
	for {
		timeBegin := time.Now()
		err := checkpointUpdater(context.Background(), UpdateCheckpointsTimeout)
		log.Err(err)
		log.Info().Interface("UpdateEpochCheckpoint took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(UpdateCheckpointsTimeout)
	}
}

func checkpointUpdater(ctx context.Context, contextTimeout time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()
	chkPoint := beacon_models.ValidatorsEpochCheckpoint{}
	err := chkPoint.GetFirstEpochCheckpointWithBalancesRemaining(ctx)
	if err != nil {
		log.Info().Err(err).Msg("checkpointUpdater")
		return err
	}
	log.Info().Msgf("UpdateEpochCheckpoint: checkpointUpdater at Epoch %d", chkPoint.Epoch)
	err = beacon_models.UpdateEpochCheckpointBalancesRecordedAtEpoch(ctxTimeout, chkPoint.Epoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalances: checkpointUpdater")
		return err
	}
	return err
}

// InsertNewEpochCheckpoint // Routine FIVE
func InsertNewEpochCheckpoint() {
	log.Info().Msg("InsertNewEpochCheckpoint")
	for {
		timeBegin := time.Now()
		err := newCheckpoint(context.Background(), InsertCheckpointsTimeout)
		log.Err(err)
		log.Info().Interface("InsertNewEpochCheckpoint took this many seconds to complete: ", time.Now().Sub(timeBegin))
		//time.Sleep(InsertCheckpointsTimeout)
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
