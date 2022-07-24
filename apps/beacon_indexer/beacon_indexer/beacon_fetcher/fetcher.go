package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var fetcher BeaconFetcher

var NewValidatorBatchSize = 1000
var NewValidatorBalancesBatchSize = 1000
var NewValidatorBalancesTimeout = time.Second * 20

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
		ctx := context.Background()
		timeBegin := time.Now()
		fetchValidatorsToInsert(ctx, NewValidatorBatchSize, sleepTime)
		log.Info().Interface("FetchNewOrMissingValidators took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(sleepTime)
	}
}

func FetchFindAndQueryAndUpdateValidatorBalances() {
	log.Info().Msg("FetchFindAndQueryAndUpdateValidatorBalances")

	for {
		ctx := context.Background()
		timeBegin := time.Now()
		fetchAndUpdateValidatorBalances(ctx, NewValidatorBalancesBatchSize, NewValidatorBalancesTimeout)
		log.Info().Interface("FetchFindAndQueryAndUpdateValidatorBalances took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(NewValidatorBalancesTimeout)
	}
}

func fetchValidatorsToInsert(ctx context.Context, batchSize int, contextTimeout time.Duration) {
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	err := fetcher.BeaconFindNewAndMissingValidatorIndexes(ctxTimeout, batchSize)
	log.Info().Err(err).Msg("FetchNewOrMissingValidators")
}

func fetchAndUpdateValidatorBalances(ctx context.Context, batchSize int, contextTimeout time.Duration) {
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	err := fetcher.FindAndQueryAndUpdateValidatorBalances(ctxTimeout, batchSize)
	log.Info().Err(err).Msg("FetchNewOrMissingValidators")
}
