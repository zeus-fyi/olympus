package iris_redis

import (
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

const (
	TwentyFive             = 25
	Fifty                  = 50
	TwoHundredFifty        = 250
	OneThousand            = 1_000
	TenThousand            = 10_000
	FiveThousand           = 5_000
	TwentyFiveThousand     = 25_000
	FiftyThousand          = 50_000
	HundredThousand        = 100_000
	FiftyMillion           = 50_000_000
	TwoHundredFiftyMillion = 250_000_000
	OneBillion             = 1_000_000_000
	ThreeBillion           = 3_000_000_000
)

func GetMonthlyPlanBudgetThroughputZU(planName string) int {
	switch strings.ToLower(planName) {
	case "enterprise":
		return HundredThousand
	case "performance":
		return HundredThousand
	case "standard":
		return FiftyThousand
	case "lite":
		return TwentyFiveThousand
	case "discover":
		return TwentyFiveThousand
	case "test":
		return 1000
	default:
		return 0
	}
}

func GetMonthlyPlanBudgetZU(planName string) int {
	switch strings.ToLower(planName) {
	case "enterprise":
		return ThreeBillion
	case "performance":
		return ThreeBillion
	case "standard":
		return OneBillion
	case "lite":
		return TwoHundredFiftyMillion
	case "discover":
		return FiftyMillion
	case "test":
		return 1000
	default:
		return 0
	}
}

func GetMonthlyPlanMaxAnvilServerlessSessions(planName string) int {
	switch strings.ToLower(planName) {
	case "enterprise":
		return MaxActiveServerlessSessions * 10
	case "performance":
		return MaxActiveServerlessSessions
	case "standard":
		return MaxActiveServerlessSessions
	case "lite":
		return MaxActiveServerlessSessions
	case "discover":
		return MaxActiveServerlessSessions
	case "test":
		return MaxActiveServerlessSessions
	default:
		return 0
	}
}

func (m *IrisCache) CheckRateLimitBroadcast(ctx context.Context, orgID int, procedureName, plan, routeGroup string, meter *iris_usage_meters.PayloadSizeMeter) (iris_programmable_proxy_v1_beta.IrisRoutingProcedure, []iris_models.RouteInfo, error) {
	// Generate the rate limiter key with the Unix timestamp
	proc, ri, um, err := m.RecordRequestUsageRatesCheckLimitAndGetBroadcastRoutes(ctx, orgID, procedureName, routeGroup, meter)
	if err != nil {
		log.Err(err).Interface("um", um).Interface("ri", ri).Msg("CheckRateLimit: RecordRequestUsageRatesCheckLimitAndNextRoute")
		return proc, ri, err
	}

	rateLimited, monthlyLimited := false, false
	switch plan {
	case "enterprise":
		// todo
		// check 100k ZU/s
		// check max 3B ZU/month
		rateLimited, monthlyLimited = um.IsRateLimited(HundredThousand, ThreeBillion)
	case "performance":
		// check 100k ZU/s
		// check max 3B ZU/month
		rateLimited, monthlyLimited = um.IsRateLimited(HundredThousand, ThreeBillion)
	case "standard":
		rateLimited, monthlyLimited = um.IsRateLimited(FiftyThousand, OneBillion)
	case "lite":
		rateLimited, monthlyLimited = um.IsRateLimited(TwentyFiveThousand, TwoHundredFiftyMillion)
	case "discover":
		rateLimited, monthlyLimited = um.IsRateLimited(FiveThousand, FiftyMillion)
	case "test":
		rateLimited, monthlyLimited = um.IsRateLimited(100, 1000)
	default:
		rateLimited, monthlyLimited = um.IsRateLimited(0, 0)
	}
	if rateLimited {
		return proc, ri, errors.New("rate limited")
	}
	if monthlyLimited {
		return proc, ri, errors.New("monthly usage exceeds plan credits")
	}
	return proc, ri, nil
}

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
		// check 100k ZU/s
		// check max 3B ZU/month
		rateLimited, monthlyLimited = um.IsRateLimited(HundredThousand, ThreeBillion)
	case "performance":
		// check 100k ZU/s
		// check max 3B ZU/month
		rateLimited, monthlyLimited = um.IsRateLimited(HundredThousand, ThreeBillion)
	case "standard":
		rateLimited, monthlyLimited = um.IsRateLimited(FiftyThousand, OneBillion)
	case "lite":
		rateLimited, monthlyLimited = um.IsRateLimited(TwentyFiveThousand, TwoHundredFiftyMillion)
	case "discover":
		rateLimited, monthlyLimited = um.IsRateLimited(FiveThousand, FiftyMillion)
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
