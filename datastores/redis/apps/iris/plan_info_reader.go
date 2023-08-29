package iris_redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

func (m *IrisCache) GetPlanUsageInfo(ctx context.Context, orgID int) (*iris_usage_meters.UsageMeter, error) {
	rateLimiterKey := getOrgRateLimitKey(orgID)
	orgRequestsMonthly := getOrgMonthlyUsageKey(orgID, time.Now().UTC().Month().String())
	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Reader.TxPipeline()

	// Get the values from Redis
	rateLimitCmd := pipe.Get(ctx, rateLimiterKey)
	monthlyUsageCmd := pipe.Get(ctx, orgRequestsMonthly)

	// Execute the pipeline
	_, err := pipe.Exec(ctx)
	if err == redis.Nil {
		err = nil
	}

	// Get the values from the commands
	rateLimit, err := rateLimitCmd.Result()
	rateLimitVal := 0.0
	if err != nil && err != redis.Nil {
		fmt.Printf("error getting rate limit: %s\n", err.Error())
		return nil, err
	} else if err == nil {
		// Convert the value to float
		rateLimitVal, err = strconv.ParseFloat(rateLimit, 64)
		if err != nil {
			fmt.Printf("error converting rate limit to float: %s\n", err.Error())
			return nil, err
		}
	}
	monthlyUsage, err := monthlyUsageCmd.Result()
	monthlyUsageVal := 0.0
	if err != nil && err != redis.Nil {
		fmt.Printf("error getting monthly usage: %s\n", err.Error())
		return nil, err
	} else if err == nil {
		// Convert the value to float
		monthlyUsageVal, err = strconv.ParseFloat(monthlyUsage, 64)
		if err != nil {
			fmt.Printf("error converting monthly usage to float: %s\n", err.Error())
			return nil, err
		}
	}
	return &iris_usage_meters.UsageMeter{
		CurrentRate:  rateLimitVal,
		MonthlyUsage: monthlyUsageVal,
	}, nil
}
