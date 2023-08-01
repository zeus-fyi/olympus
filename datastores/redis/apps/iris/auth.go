package iris_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	util "github.com/wealdtech/go-eth2-util"
)

func (m *IrisCache) GetAuthCacheIfExists(ctx context.Context, token string) (int64, string, error) {
	// Generate the rate limiter key with the Unix timestamp
	hashedToken := fmt.Sprintf("%x", util.Keccak256([]byte(token)))
	hashedTokenPlan := fmt.Sprintf("%x:plan", util.Keccak256([]byte(token)))

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

func (m *IrisCache) SetAuthCache(ctx context.Context, orgID int, token, plan string) error {
	// Generate the rate limiter key with the Unix timestamp
	hashedToken := fmt.Sprintf("%x", util.Keccak256([]byte(token)))
	hashedTokenPlan := fmt.Sprintf("%x:plan", util.Keccak256([]byte(token)))

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()
	pipe.Set(ctx, hashedToken, orgID, time.Minute*15)
	pipe.Set(ctx, hashedTokenPlan, plan, time.Minute*15)
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("SetAuthCache")
		return err
	}
	return nil
}
