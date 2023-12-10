package iris_redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
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

func (r *IrisRedisTestSuite) TestGetAllOrgMonthlyUsage() {
	meter := &iris_usage_meters.PayloadSizeMeter{}
	meter.Month = time.Now().UTC().Month().String()
	ot, aerr := iris_models.SelectAllOrgRoutes(context.Background())
	r.Require().NoError(aerr)
	for orgID, _ := range ot.Map {
		orgRequestsMonthly := getOrgMonthlyUsageKey(orgID, meter.Month)
		// Use Redis transaction (pipeline) to perform all operations atomically
		pipe := IrisRedisClient.Reader.TxPipeline()
		// Increment the payload size meter
		ctx := context.Background()

		rateLimitCmd := pipe.Get(ctx, orgRequestsMonthly)

		// Execute the pipeline
		_, err := pipe.Exec(ctx)
		if err == redis.Nil {
			//fmt.Println("key doesn't exist: ", orgRequestsMonthly)
			continue
		}
		// Get the values from the commands
		rateLimit, err := rateLimitCmd.Result()
		r.Assert().NoError(err)
		// Convert the value to float
		rateLimitVal, err := strconv.ParseFloat(rateLimit, 64)
		r.Assert().NoError(err)
		fmt.Println(orgRequestsMonthly, rateLimitVal)
	}
}
