package beacon_fetcher

import (
	"github.com/rs/zerolog"
)

var fetcher BeaconFetcher

func InitFetcherService(nodeURL string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	fetcher.NodeEndpoint = nodeURL
	go FetchNewOrMissingValidators()
	go FetchAllValidatorBalances()
	go UpdateAllValidators()
	go UpdateEpochCheckpoint()
	go InsertNewEpochCheckpoint()
}
