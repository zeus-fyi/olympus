package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var fetcher BeaconFetcher

func InitFetcherService(nodeURL string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	fetcher.NodeEndpoint = nodeURL
	fetchNewValidatorTimeout := time.Minute * 5
	go FetchNewOrMissingValidators(10000, fetchNewValidatorTimeout)
	fetchUpdateTimeout := time.Second * 5
	go FetchFindAndQueryAndUpdateValidatorBalances(200000, fetchUpdateTimeout)
}

func FetchNewOrMissingValidators(batchSize int, sleepTime time.Duration) {
	log.Info().Msg("FetchNewOrMissingValidators")

	for {
		ctx := context.Background()
		fetchValidatorsToInsert(ctx, batchSize, sleepTime)
		time.Sleep(sleepTime)
	}
}

func fetchValidatorsToInsert(ctx context.Context, batchSize int, contextTimeout time.Duration) {
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	err := fetcher.BeaconFindNewAndMissingValidatorIndexes(ctxTimeout, batchSize)
	log.Info().Err(err).Msg("FetchNewOrMissingValidators")
}

func FetchFindAndQueryAndUpdateValidatorBalances(batchSize int, sleepTime time.Duration) {
	log.Info().Msg("FetchFindAndQueryAndUpdateValidatorBalances")

	for {
		ctx := context.Background()
		fetchAndUpdateValidatorBalances(ctx, batchSize, sleepTime)
		time.Sleep(sleepTime)

	}
}

func fetchAndUpdateValidatorBalances(ctx context.Context, batchSize int, contextTimeout time.Duration) {
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	err := fetcher.FindAndQueryAndUpdateValidatorBalances(ctxTimeout, batchSize)
	log.Info().Err(err).Msg("FetchNewOrMissingValidators")
}
