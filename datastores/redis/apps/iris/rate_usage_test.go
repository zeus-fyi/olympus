package iris_redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

func (r *IrisRedisTestSuite) TestNewRateUsage() {
	meter := &iris_usage_meters.PayloadSizeMeter{}
	orgID := 420
	rateLimiterKey := getOrgRateLimitKey(orgID)
	orgRequestsTotal := getOrgTotalRequestsKey(orgID)
	if len(meter.Month) <= 0 {
		meter.Month = time.Now().UTC().Month().String()
	}
	orgRequestsMonthly := getOrgMonthlyUsageKey(orgID, meter.Month)
	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := IrisRedisClient.Writer.TxPipeline()
	// Increment the payload size meter
	ctx := context.Background()
	_ = pipe.IncrByFloat(ctx, orgRequestsTotal, meter.ZeusRequestComputeUnitsConsumed())
	_ = pipe.IncrByFloat(ctx, orgRequestsMonthly, meter.ZeusRequestComputeUnitsConsumed())
	// Increment the rate limiter key
	_ = pipe.IncrByFloat(ctx, rateLimiterKey, meter.ZeusRequestComputeUnitsConsumed())
	// rate limiter key expires after 3 seconds
	pipe.Expire(ctx, rateLimiterKey, time.Second*3)
	// Execute the pipeline
	_, err := pipe.Exec(ctx)
	if err == redis.Nil {
		err = nil
	}
}
