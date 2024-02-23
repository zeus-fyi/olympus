package beacon_fetcher

import (
	"context"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog"
	"github.com/zeus-fyi/olympus/datastores/redis/apps/beacon_indexer"
)

var (
	Fetcher   BeaconFetcher
	NetworkID = 1
)

func InitFetcherService(ctx context.Context, nodeURL string, redis *redis.Client) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	Fetcher.NodeEndpoint = nodeURL
	Fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, redis)
	minv := 1
	maxv := 20
	jitterStart := time.Duration(rand.Intn(maxv-minv+1) + minv)
	time.Sleep(time.Second * jitterStart)
	go FetchNewOrMissingValidators(NetworkID)
	go UpdateAllValidatorBalancesFromCache()
	// used for redis look ahead cache
	go FetchAnyValidatorBalancesAfterCheckpoint()
	go FetchBeaconUpdateValidatorStates(NetworkID)
	go UpdateAllValidators()
	go UpdateForwardEpochCheckpoint()
	go InsertNewEpochCheckpoint()
}
