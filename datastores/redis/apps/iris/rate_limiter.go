package iris_redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v9"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

const (
	OneThousand            = 1_000
	TenThousand            = 10_000
	TwentyFiveThousand     = 25_000
	FiftyThousand          = 50_000
	FiftyMillion           = 50_000_000
	TwoHundredFiftyMillion = 250_000_000
	OneBillion             = 1_000_000_000
	ThreeBillion           = 3_000_000_000
)

func (m *IrisCache) CheckRateLimit(ctx context.Context, orgID int, plan string, meter *iris_usage_meters.PayloadSizeMeter) error {
	// Generate the rate limiter key with the Unix timestamp
	um, err := m.GetUsageRates(ctx, orgID, meter)
	if err != nil {
		return err
	}
	rateLimited, monthlyLimited := false, false
	switch plan {
	case "performance":
		// check 50k ZU/s
		// check max 3B ZU/month
		rateLimited, monthlyLimited = um.IsRateLimited(FiftyThousand, ThreeBillion)
	case "standard":
		// check 25k ZU/s
		// check max 1B ZU/month
		rateLimited, monthlyLimited = um.IsRateLimited(TwentyFiveThousand, OneBillion)
	case "lite":
		// check 1k ZU/s
		// check max 50M ZU/month
		rateLimited, monthlyLimited = um.IsRateLimited(TenThousand, TwoHundredFiftyMillion)
	case "test":
		// check 1k ZU/s
		// check max 50M ZU/month
		rateLimited, monthlyLimited = um.IsRateLimited(100, 1000)
	default:
		rateLimited, monthlyLimited = um.IsRateLimited(0, 0)
	}
	if rateLimited {
		return errors.New("rate limited")
	}
	if monthlyLimited {
		return errors.New("monthly usage exceeds plan credits")
	}
	return nil
}

func (m *IrisCache) GetUsageRates(ctx context.Context, orgID int, meter *iris_usage_meters.PayloadSizeMeter) (iris_usage_meters.UsageMeter, error) {
	// Generate the rate limiter key with the Unix timestamp
	rateLimiterKey := orgRateLimitTag(orgID)
	orgRequests := fmt.Sprintf("%d-total-zu-usage", orgID)
	if meter != nil {
		orgRequests = orgMonthlyUsageTag(orgID, meter.Month)
	}

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Reader.TxPipeline()

	// Get the values from Redis
	rateLimitCmd := pipe.Get(ctx, rateLimiterKey)
	monthlyUsageCmd := pipe.Get(ctx, orgRequests)

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
		return iris_usage_meters.UsageMeter{}, err
	} else if err == nil {
		// Convert the value to float
		rateLimitVal, err = strconv.ParseFloat(rateLimit, 64)
		if err != nil {
			fmt.Printf("error converting rate limit to float: %s\n", err.Error())
			return iris_usage_meters.UsageMeter{}, err
		}
	}

	monthlyUsage, err := monthlyUsageCmd.Result()
	monthlyUsageVal := 0.0
	if err != nil && err != redis.Nil {
		fmt.Printf("error getting monthly usage: %s\n", err.Error())
		return iris_usage_meters.UsageMeter{}, err
	} else if err == nil {
		// Convert the value to float
		monthlyUsageVal, err = strconv.ParseFloat(monthlyUsage, 64)
		if err != nil {
			fmt.Printf("error converting monthly usage to float: %s\n", err.Error())
			return iris_usage_meters.UsageMeter{}, err
		}
	}

	// Return the UsageRates
	return iris_usage_meters.UsageMeter{
		RateLimit:    rateLimitVal,
		MonthlyUsage: monthlyUsageVal,
	}, nil
}

/*
var LocalRateLimiterCache = cache.New(1*time.Second, 2*time.Second)
_, found := LocalRateLimiterCache.Get(rateLimiterKey)

	if found {
		err := LocalRateLimiterCache.Increment(rateLimiterKey, 1)
		if err != nil {
			log.Err(err).Msg("LocalRateLimiterCache: GetNextRoute")
		}
	} else {

		LocalRateLimiterCache.Set(rateLimiterKey, 1, 1*time.Second)
	}
*/
