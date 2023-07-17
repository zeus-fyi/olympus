package redis_mev

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

// todo v3

func (m *MevCache) AddOrUpdateV2PairPricesCache(ctx context.Context, blockNumber uint64, ttl time.Duration) error {
	statusCmd := m.Set(ctx, LatestBlockNumberCacheKey, blockNumber, ttl)
	if statusCmd.Err() != nil {
		log.Ctx(ctx).Err(statusCmd.Err()).Msgf("AddOrUpdateLatestBlockCache: %d", blockNumber)
		return statusCmd.Err()
	}
	return nil
}

func (m *MevCache) GetV2PairPricesFromCacheIfExists(ctx context.Context) (uint64, error) {
	val, err := m.Get(ctx, LatestBlockNumberCacheKey).Uint64()
	switch {
	case err == redis.Nil:
		log.Info().Msgf("GetLatestV2PairPrices: latest v2 pair not in cache")
		return 0, errors.New("latest block number not in cache")
	case err != nil:
		log.Err(err).Msgf("GetLatestBlockNumber")
		return 0, err
	}
	return val, err
}
