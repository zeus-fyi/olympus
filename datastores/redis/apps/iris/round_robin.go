package iris_redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

func orgRouteTag(orgID int, rgName string) string {
	return fmt.Sprintf("%d-%s", orgID, rgName)
}

func orgMonthlyUsageTag(orgID int, month string) string {
	return fmt.Sprintf("%d-%s-total-zu-usage", orgID, month)
}

func orgRateLimitTag(orgID int) string {
	return fmt.Sprintf("%d:%d", orgID, time.Now().Unix())
}

func (m *IrisCache) GetNextRoute(ctx context.Context, orgID int, rgName string, meter *iris_usage_meters.PayloadSizeMeter) (iris_models.RouteInfo, error) {
	// Generate the rate limiter key with the Unix timestamp
	rateLimiterKey := orgRateLimitTag(orgID)
	orgRequests := fmt.Sprintf("%d-total-zu-usage", orgID)
	if meter != nil {
		orgRequests = orgMonthlyUsageTag(orgID, meter.Month)
	}

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()

	if meter != nil {
		// Increment the payload size meter
		_ = pipe.IncrByFloat(ctx, orgRequests, meter.ZeusRequestComputeUnitsConsumed())
		// Increment the rate limiter key
		_ = pipe.IncrByFloat(ctx, rateLimiterKey, meter.ZeusRequestComputeUnitsConsumed())
	} else {
		_ = pipe.Incr(ctx, orgRequests)
		// Increment the rate limiter key
		_ = pipe.Incr(ctx, rateLimiterKey)
	}

	// Set the key to expire after 2 seconds
	pipe.Expire(ctx, rateLimiterKey, 10*time.Second)

	// Generate the route key
	routeKey := orgRouteTag(orgID, rgName)

	// Check if the key exists
	existsCmd := pipe.Exists(ctx, routeKey)

	// Pop the endpoint from the head of the list and push it back to the tail, ensuring round-robin rotation
	endpointCmd := pipe.LMove(ctx, routeKey, routeKey, "LEFT", "RIGHT")

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		// If there's an error, log it and return it
		fmt.Printf("error during pipeline execution: %s\n", err.Error())
		return iris_models.RouteInfo{}, err
	}

	// Check whether the key exists
	if existsCmd.Val() <= 0 {
		log.Err(err).Msgf("key doesn't exist: %s", routeKey)
		return iris_models.RouteInfo{}, fmt.Errorf("key doesn't exist: %s", routeKey)
	}

	// Get referers from Redis
	refererKey := orgRouteTag(orgID, endpointCmd.Val()) + ":referers"
	referers, err := m.Reader.SMembers(ctx, refererKey).Result()
	if err != nil {
		log.Err(err).Msgf("error getting referers for route: %s", refererKey)
		return iris_models.RouteInfo{}, fmt.Errorf("error getting referers for route: %s", refererKey)
	}

	// Return the RouteInfo
	return iris_models.RouteInfo{
		RoutePath: endpointCmd.Val(),
		Referers:  referers,
	}, nil
}

func (m *IrisCache) IncrementResponseUsageRateMeter(ctx context.Context, orgID int, meter *iris_usage_meters.PayloadSizeMeter) error {
	// Generate the rate limiter key with the Unix timestamp
	rateLimiterKey := orgRateLimitTag(orgID)
	orgRequests := fmt.Sprintf("%d-total-zu-usage", orgID)
	if meter != nil {
		orgRequests = orgMonthlyUsageTag(orgID, meter.Month)
	}

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()
	if meter != nil {
		// Increment the payload size meter
		_ = pipe.IncrByFloat(ctx, orgRequests, meter.ZeusResponseComputeUnitsConsumed())
		// Increment the rate limiter key
		_ = pipe.IncrByFloat(ctx, rateLimiterKey, meter.ZeusResponseComputeUnitsConsumed())
	} else {
		_ = pipe.Incr(ctx, orgRequests)
		// Increment the rate limiter key
		_ = pipe.Incr(ctx, rateLimiterKey)
	}
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("AddOrUpdateOrgRoutingGroup")
		return err
	}
	return nil
}

func (m *IrisCache) AddOrUpdateOrgRoutingGroup(ctx context.Context, orgID int, rgName string, routes []iris_models.RouteInfo) error {
	// Start a new transaction
	pipe := m.Writer.TxPipeline()

	// Define the route group tag
	rgTag := orgRouteTag(orgID, rgName)

	// Remove the old route group key if it exists
	pipe.Del(ctx, rgTag)

	// Add each route to the list
	for _, routeInfo := range routes {
		// Define unique tags for route and referers
		routeTag := orgRouteTag(orgID, routeInfo.RoutePath)
		refererTag := routeTag + ":referers"

		// Remove the old keys if they exist
		pipe.Del(ctx, routeTag, refererTag)

		// Add the route
		pipe.RPush(ctx, rgTag, routeInfo.RoutePath)

		// Add the referers to a set
		if len(routeInfo.Referers) > 0 {
			// Convert []string to []interface{}
			referersInterface := make([]interface{}, len(routeInfo.Referers))
			for i, v := range routeInfo.Referers {
				referersInterface[i] = v
			}
			pipe.SAdd(ctx, refererTag, referersInterface...)
		}
	}

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("AddOrUpdateOrgRoutingGroup: %s", rgTag)
		return err
	}
	return nil
}

func (m *IrisCache) DeleteOrgRoutingGroup(ctx context.Context, orgID int, rgName string) error {
	// Generate the key
	tag := orgRouteTag(orgID, rgName)

	// Start a new transaction
	pipe := m.Writer.TxPipeline()

	// Remove the old key if it exists
	pipe.Del(ctx, tag)

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("DeleteOrgRoutingGroup: %s", tag)
		return err
	}
	return nil
}

func (m *IrisCache) initRoutingTables(ctx context.Context) error {
	ot, err := iris_models.SelectAllOrgRoutes(ctx)
	if err != nil {
		return err
	}
	for orgID, og := range ot.Map {
		for rgName, routes := range og {
			err = m.AddOrUpdateOrgRoutingGroup(context.Background(), orgID, rgName, routes)
			if err != nil {
				log.Err(err).Msg("InitRoutingTables")
				return err
			}
		}
	}
	return nil
}

func (m *IrisCache) initRoutingTablesForOrg(ctx context.Context, orgID int) error {
	ot, err := iris_models.SelectAllOrgRoutesByOrg(ctx, orgID)
	if err != nil {
		return err
	}
	for _, og := range ot.Map {
		for rgName, routes := range og {
			err = m.AddOrUpdateOrgRoutingGroup(context.Background(), orgID, rgName, routes)
			if err != nil {
				log.Err(err).Msg("InitRoutingTablesForOrg")
				return err
			}
		}
	}
	return nil
}

func (m *IrisCache) refreshRoutingTablesForOrgGroup(ctx context.Context, orgID int, groupName string) error {
	ot, err := iris_models.SelectOrgRoutesByOrgAndGroupName(ctx, orgID, groupName)
	if err != nil {
		return err
	}
	for _, og := range ot.Map {
		for rgName, routes := range og {
			err = m.AddOrUpdateOrgRoutingGroup(context.Background(), orgID, rgName, routes)
			if err != nil {
				log.Err(err).Msg("InitRoutingTablesForOrgGroup")
				return err
			}
		}
	}
	return nil
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
const (
	OneThousand        = 1_000
	TwentyFiveThousand = 25_000
	FiftyThousand      = 50_000
	FiftyMillion       = 50_000_000
	OneBillion         = 1_000_000_000
	ThreeBillion       = 3_000_000_000
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
	case "free":
		// check 1k ZU/s
		// check max 50M ZU/month
		rateLimited, monthlyLimited = um.IsRateLimited(OneThousand, FiftyMillion)
	case "test":
		// check 1k ZU/s
		// check max 50M ZU/month
		rateLimited, monthlyLimited = um.IsRateLimited(100, 1000)
	default:
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
