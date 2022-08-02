package beacon_indexer

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

type FetcherCache struct {
	*redis.Client
}

func NewFetcherCache(ctx context.Context, r *redis.Client) FetcherCache {
	log.Ctx(ctx).Info().Msg("NewFetcherCache")
	log.Info().Interface("redis", r)
	return FetcherCache{r}
}

func (f *FetcherCache) SetCheckpointCache(ctx context.Context, epoch int, ttl time.Duration) (string, error) {
	key := fmt.Sprintf("checkpoint-epoch-%d", epoch)
	statusCmd := f.Set(ctx, fmt.Sprintf("checkpoint-epoch-%d", epoch), epoch, ttl)
	if statusCmd.Err() != nil {
		log.Ctx(ctx).Err(statusCmd.Err())
		return "", statusCmd.Err()
	}
	log.Ctx(ctx).Info().Msgf("set cache at epoch %d", epoch)
	return key, nil
}

func (f *FetcherCache) DoesCheckpointExist(ctx context.Context, epoch int) bool {
	key := fmt.Sprintf("checkpoint-epoch-%d", epoch)
	chkPoint, err := f.Get(ctx, key).Int()
	if err != nil {
		return false
	}

	return chkPoint == epoch
}

func (f *FetcherCache) DeleteCheckpoint(ctx context.Context, epoch int) error {
	key := fmt.Sprintf("checkpoint-epoch-%d", epoch)
	err := f.Del(ctx, key)
	if err != nil {
		return err.Err()
	}
	return err.Err()
}
