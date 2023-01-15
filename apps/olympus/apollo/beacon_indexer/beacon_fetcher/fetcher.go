package beacon_fetcher

import (
	"context"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog"
	"github.com/zeus-fyi/olympus/datastores/redis/apps/beacon_indexer"
)

var Fetcher BeaconFetcher

func InitFetcherService(ctx context.Context, nodeURL string, redis *redis.Client) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	Fetcher.NodeEndpoint = nodeURL
	Fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, redis)
	min := 1
	max := 20
	jitterStart := time.Duration(rand.Intn(max-min+1) + min)
	time.Sleep(time.Second * jitterStart)
	go FetchNewOrMissingValidators()
	go UpdateAllValidatorBalancesFromCache()
	// used for redis look ahead cache
	go FetchAnyValidatorBalancesAfterCheckpoint()
	go FetchBeaconUpdateValidatorStates(1)
	go UpdateAllValidators()
	go UpdateForwardEpochCheckpoint()
	go InsertNewEpochCheckpoint()
}
