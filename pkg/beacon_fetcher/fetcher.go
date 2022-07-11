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

	go FetchNewOrMissingValidators()
	go FetchFindAndQueryAndUpdateValidatorBalances()
}

func FetchNewOrMissingValidators() {
	log.Info().Msg("FetchNewOrMissingValidators")

	sleepBetweenFetches := time.Minute * 5
	batchSize := 1000
	for {
		ctx := context.Background()
		err := fetcher.BeaconFindNewAndMissingValidatorIndexes(ctx, batchSize)
		log.Info().Err(err).Msg("FetchNewOrMissingValidators")
		time.Sleep(sleepBetweenFetches)
	}
}

func FetchFindAndQueryAndUpdateValidatorBalances() {
	log.Info().Msg("FetchFindAndQueryAndUpdateValidatorBalances")

	sleepBetweenFetches := time.Second * 20
	batchSize := 10000
	for {
		ctx := context.Background()
		err := fetcher.FindAndQueryAndUpdateValidatorBalances(ctx, batchSize)
		log.Info().Err(err).Msg("FetchNewOrMissingValidators")
		time.Sleep(sleepBetweenFetches)
	}
}
