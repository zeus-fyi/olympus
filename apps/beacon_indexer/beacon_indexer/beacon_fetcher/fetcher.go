package beacon_fetcher

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog"
	"github.com/zeus-fyi/olympus/pkg/datastores/redis_app"
	"github.com/zeus-fyi/olympus/pkg/datastores/redis_app/beacon_indexer"
)

var fetcher BeaconFetcher

func InitFetcherService(nodeURL string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	ctx := context.Background()

	fetcher.NodeEndpoint = nodeURL

	redisOpts := redis.Options{
		Addr: "localhost:6379",
	}
	r := redis_app.InitRedis(ctx, redisOpts)
	fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, r)

	go FetchNewOrMissingValidators()
	// go FetchAllValidatorBalances()
	go FetchAllValidatorBalancesAfterCheckpoint()
	go UpdateAllValidators()
	go UpdateEpochCheckpoint()
	go InsertNewEpochCheckpoint()
}
