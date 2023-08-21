package iris_redis

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
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

func (m *IrisCache) CheckRateLimit(ctx context.Context, orgID int, plan, routeGroup string, meter *iris_usage_meters.PayloadSizeMeter) (iris_models.RouteInfo, error) {
	// Generate the rate limiter key with the Unix timestamp
	ri, um, err := m.RecordRequestUsageRatesCheckLimitAndNextRoute(ctx, orgID, routeGroup, meter)
	if err != nil {
		log.Err(err).Interface("um", um).Interface("ri", ri).Msg("CheckRateLimit: RecordRequestUsageRatesCheckLimitAndNextRoute")
		return ri, err
	}
	rateLimited, monthlyLimited := false, false
	switch plan {
	case "enterprise":
		// todo
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
		return ri, errors.New("rate limited")
	}
	if monthlyLimited {
		return ri, errors.New("monthly usage exceeds plan credits")
	}
	return ri, nil
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
