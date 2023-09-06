package iris_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

func (m *IrisCache) GetAuthCacheIfExists(ctx context.Context, token string) (int64, string, error) {
	// Generate the rate limiter key with the Unix timestamp
	hashedToken := getHashedTokenKey(token)
	hashedTokenPlan := getHashedTokenPlanKey(token)

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Reader.TxPipeline()

	// Get the values from Redis
	orgIDCmd := pipe.Get(ctx, hashedToken)
	orgPlanCmd := pipe.Get(ctx, hashedTokenPlan)

	// Execute the pipeline
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return -1, "", err
	}

	// Get the values from the commands
	orgID, err := orgIDCmd.Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, "", fmt.Errorf("no value found for token: %s", token)
		}
		return -1, "", err
	}

	plan, err := orgPlanCmd.Result()
	if err != nil {
		if err == redis.Nil {
			return 0, "", fmt.Errorf("no value found for token: %s", token)
		}
		return -1, "", err
	}

	return orgID, plan, nil
}

const (
	CacheApiKeyTTL  = time.Minute * 15
	CacheSessionTTL = time.Second * 3
)

func (m *IrisCache) SetAuthCache(ctx context.Context, orgID int, token, plan string, usingCookie bool) error {
	// Generate the rate limiter key with the Unix timestamp
	hashedToken := getHashedTokenKey(token)
	hashedTokenPlan := getHashedTokenPlanKey(token)

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()
	ttl := CacheApiKeyTTL
	if usingCookie {
		ttl = CacheSessionTTL
	}
	pipe.Set(ctx, hashedToken, orgID, ttl)
	pipe.Set(ctx, hashedTokenPlan, plan, ttl)
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("SetAuthCache")
		return err
	}
	return nil
}

func (m *IrisCache) DeleteAuthCache(ctx context.Context, token string) error {
	// Generate the rate limiter key with the Unix timestamp
	hashedToken := getHashedTokenKey(token)
	hashedTokenPlan := getHashedTokenPlanKey(token)

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()
	pipe.Del(ctx, hashedToken)
	pipe.Del(ctx, hashedTokenPlan)
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("DeleteAuthCache")
		return err
	}
	return nil
}
