package beacon_fetcher

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog"
	"github.com/zeus-fyi/olympus/pkg/datastores/redis_app/beacon_indexer"
)

var fetcher BeaconFetcher

func InitFetcherService(ctx context.Context, nodeURL string, redis *redis.Client) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	fetcher.NodeEndpoint = nodeURL
	fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, redis)

	go FetchNewOrMissingValidators()
	// go FetchAllValidatorBalances()
	go FetchAllValidatorBalancesAfterCheckpoint()
	go UpdateAllValidators()
	//go UpdateEpochCheckpoint()
	go UpdateForwardEpochCheckpoint()
	go InsertNewEpochCheckpoint()
}
