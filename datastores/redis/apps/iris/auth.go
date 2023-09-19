package iris_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (m *IrisCache) GetAuthCacheIfExists(ctx context.Context, token string) (org_users.OrgUser, string, error) {
	// Generate the rate limiter key with the Unix timestamp
	hashedToken := getHashedTokenKey(token)
	hashedTokenPlan := getHashedTokenPlanKey(token)
	hashedTokenUserID := getHashedTokenUserID(token)

	ou := org_users.OrgUser{
		OrgUsers: autogen_bases.OrgUsers{UserID: int(-1), OrgID: int(-1)},
	}

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Reader.TxPipeline()

	// Get the values from Redis
	orgIDCmd := pipe.Get(ctx, hashedToken)
	orgPlanCmd := pipe.Get(ctx, hashedTokenPlan)
	orgUserIDCmd := pipe.Get(ctx, hashedTokenUserID)
	// Execute the pipeline
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return ou, "", err
	}

	// Get the values from the commands
	orgID, err := orgIDCmd.Int64()
	if err != nil {
		if err == redis.Nil {
			return ou, "", fmt.Errorf("no value found for hashedToken: %s", hashedToken)
		}
		return ou, "", err
	}

	plan, err := orgPlanCmd.Result()
	if err != nil {
		if err == redis.Nil {
			return ou, "", fmt.Errorf("no value found for hashedTokenPlan: %s", hashedTokenPlan)
		}
		return ou, "", err
	}

	userID, err := orgUserIDCmd.Int64()
	if err != nil {
		if err == redis.Nil {
			return ou, "", fmt.Errorf("no value found for hashedTokenUserID: %s", hashedTokenUserID)
		}
		return ou, "", err
	}

	ou = org_users.OrgUser{
		OrgUsers: autogen_bases.OrgUsers{UserID: int(userID), OrgID: int(orgID)},
	}
	return ou, plan, nil
}

const (
	CacheApiKeyTTL  = time.Minute * 15
	CacheSessionTTL = time.Second * 3
)

func (m *IrisCache) SetAuthCache(ctx context.Context, ou org_users.OrgUser, token, plan string, usingCookie bool) error {
	// Generate the rate limiter key with the Unix timestamp
	hashedToken := getHashedTokenKey(token)
	hashedTokenPlan := getHashedTokenPlanKey(token)
	hashedTokenUserID := getHashedTokenUserID(token)

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()
	ttl := CacheApiKeyTTL
	if usingCookie {
		ttl = CacheSessionTTL
	}
	pipe.Set(ctx, hashedToken, ou.OrgID, ttl)
	pipe.Set(ctx, hashedTokenPlan, plan, ttl)
	pipe.Set(ctx, hashedTokenUserID, ou.UserID, ttl)
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
	hashedTokenUserID := getHashedTokenUserID(token)

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()
	pipe.Del(ctx, hashedToken)
	pipe.Del(ctx, hashedTokenPlan)
	pipe.Del(ctx, hashedTokenUserID)
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("DeleteAuthCache")
		return err
	}
	return nil
}
