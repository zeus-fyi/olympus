package beacon_indexer

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

type FetcherCache struct {
	*redis.Client
}

func (f *FetcherCache) CheckpointCache(ctx context.Context) {

	statusCmd := f.Set(ctx, "k", "vb", time.Second)
	log.Ctx(ctx).Info().Msgf("statusCmd %s", statusCmd)

}
