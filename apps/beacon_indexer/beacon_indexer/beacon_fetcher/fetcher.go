package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres/beacon_indexer/beacon_models"
)

var fetcher BeaconFetcher

var NewValidatorBatchSize = 100
var NewValidatorBalancesBatchSize = 1000
var NewValidatorBalancesTimeout = time.Second * 180

var NewValidatorTimeout = 60 * time.Minute
var UpdateValidatorTimeout = time.Minute * 10

var FetchAllValidatorBalancesTimeout = time.Minute * 10

// Checkpoints

var UpdateCheckpointsTimeout = 60 * time.Second
var InsertCheckpointsTimeout = time.Minute * 2

func InitFetcherService(nodeURL string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	fetcher.NodeEndpoint = nodeURL
	go FetchNewOrMissingValidators(NewValidatorTimeout)
	go FetchAllValidatorBalances()
	go UpdateAllValidators()
	go UpdateEpochCheckpoint()
	go InsertNewEpochCheckpoint()
}

// FetchNewOrMissingValidators Routine ONE
func FetchNewOrMissingValidators(sleepTime time.Duration) {
	log.Info().Msg("FetchNewOrMissingValidators")

	for {
		timeBegin := time.Now()
		err := fetchValidatorsToInsert(context.Background(), NewValidatorBatchSize, sleepTime)
		log.Err(err)
		log.Info().Interface("FetchNewOrMissingValidators took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(NewValidatorTimeout)
	}
}

func fetchValidatorsToInsert(ctx context.Context, batchSize int, contextTimeout time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()
	err := fetcher.BeaconFindNewAndMissingValidatorIndexes(ctxTimeout, batchSize)
	log.Info().Err(err).Msg("FetchNewOrMissingValidators")
	return err
}

// UpdateAllValidators Routine TWO
func UpdateAllValidators() {
	log.Info().Msg("UpdateAllValidators")
	for {
		timeBegin := time.Now()
		err := fetchAllValidatorsToUpdate(context.Background(), UpdateValidatorTimeout)
		log.Err(err)
		log.Info().Interface("UpdateAllValidators took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(UpdateValidatorTimeout)
	}
}

func fetchAllValidatorsToUpdate(ctx context.Context, contextTimeout time.Duration) error {
	log.Info().Msg("fetchAllValidatorsToUpdate")
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	err := fetcher.BeaconUpdateAllValidatorStates(ctxTimeout)
	log.Info().Err(err).Msg("UpdateAllValidators: fetchAllValidatorsToUpdate")
	return err
}

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

// Routine SIX
//func FetchFindAndQueryAndUpdateValidatorBalances() {
//	log.Info().Msg("FetchFindAndQueryAndUpdateValidatorBalances")
//
//	for {
//		timeBegin := time.Now()
//		err := fetchAndUpdateValidatorBalances(context.Background(), NewValidatorBalancesBatchSize, NewValidatorBalancesTimeout)
//		log.Err(err)
//		log.Info().Interface("FetchFindAndQueryAndUpdateValidatorBalances took this many seconds to complete: ", time.Now().Sub(timeBegin))
//	}
//}

//func fetchAndUpdateValidatorBalances(ctx context.Context, batchSize int, contextTimeout time.Duration) error {
//	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
//	defer cancel()
//
//	err := fetcher.FindAndQueryAndUpdateValidatorBalances(ctxTimeout, batchSize)
//	log.Info().Err(err).Msg("FetchNewOrMissingValidators")
//	return err
//}
