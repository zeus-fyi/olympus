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
var NewAllValidatorBalancesTimeout = time.Minute * 5

var NewValidatorTimeout = time.Minute * 60
var UpdateValidatorTimeout = time.Minute * 10

func InitFetcherService(nodeURL string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	fetcher.NodeEndpoint = nodeURL
	go FetchNewOrMissingValidators(NewValidatorTimeout)
	go FetchAllValidatorBalances()
	go UpdateAllValidators()
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
		err := fetchAllValidatorBalances(context.Background(), NewAllValidatorBalancesTimeout)
		log.Err(err)
		log.Info().Interface("FetchFindAndQueryAndUpdateValidatorBalances took this many seconds to complete: ", time.Now().Sub(timeBegin))
	}
}

func fetchAllValidatorBalances(ctx context.Context, contextTimeout time.Duration) error {
	log.Info().Msg("fetchAllValidatorBalances")
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	chkPoint := beacon_models.ValidatorsEpochCheckpoint{}
	err := beacon_models.UpdateEpochCheckpointBalancesRecordedAtEpoch(ctx, chkPoint.Epoch)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalances: UpdateEpochCheckpointBalancesRecordedAtEpoch")
		return err
	}
	err = chkPoint.GetFirstEpochCheckpointWithBalancesRemaining(ctx)
	if err != nil {
		log.Info().Err(err).Msg("fetchAllValidatorBalances")
		return err
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

	err = beacon_models.InsertEpochCheckpoint(ctx, chkPoint.Epoch)
	if err != nil {
		log.Error().Err(err).Msg("fetchAllValidatorBalances: InsertEpochCheckpoint")
		return err
	}

	log.Info().Err(err).Msg("fetchAllValidatorBalances")
	return err
}

// Routine FOUR
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
