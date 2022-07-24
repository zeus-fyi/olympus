package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var fetcher BeaconFetcher

var NewValidatorBatchSize = 1000
var NewValidatorBalancesBatchSize = 10000
var NewValidatorBalancesTimeout = time.Second * 180

func InitFetcherService(nodeURL string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	fetcher.NodeEndpoint = nodeURL
	fetchNewValidatorTimeout := time.Minute * 5
	go FetchNewOrMissingValidators(fetchNewValidatorTimeout)
	go FetchFindAndQueryAndUpdateValidatorBalances()
}

func FetchNewOrMissingValidators(sleepTime time.Duration) {
	log.Info().Msg("FetchNewOrMissingValidators")

	for {
		timeBegin := time.Now()
		err := fetchValidatorsToInsert(context.Background(), NewValidatorBatchSize, sleepTime)
		log.Err(err)
		log.Info().Interface("FetchNewOrMissingValidators took this many seconds to complete: ", time.Now().Sub(timeBegin))
	}
}

func FetchFindAndQueryAndUpdateValidatorBalances() {
	log.Info().Msg("FetchFindAndQueryAndUpdateValidatorBalances")

	for {
		timeBegin := time.Now()
		err := fetchAndUpdateValidatorBalances(context.Background(), NewValidatorBalancesBatchSize, NewValidatorBalancesTimeout)
		log.Err(err)
		log.Info().Interface("FetchFindAndQueryAndUpdateValidatorBalances took this many seconds to complete: ", time.Now().Sub(timeBegin))
	}
}

func fetchValidatorsToInsert(ctx context.Context, batchSize int, contextTimeout time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()
	err := fetcher.BeaconFindNewAndMissingValidatorIndexes(ctxTimeout, batchSize)
	log.Info().Err(err).Msg("FetchNewOrMissingValidators")
	return err
}

func fetchAndUpdateValidatorBalances(ctx context.Context, batchSize int, contextTimeout time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	err := fetcher.FindAndQueryAndUpdateValidatorBalances(ctxTimeout, batchSize)
	log.Info().Err(err).Msg("FetchNewOrMissingValidators")
	return err
}
