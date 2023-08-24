package iris_redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

func (m *IrisCache) RecordRequestUsageRatesCheckLimitAndNextRoute(ctx context.Context, orgID int, rgName string, meter *iris_usage_meters.PayloadSizeMeter) (iris_models.RouteInfo, iris_usage_meters.UsageMeter, error) {
	// Generate the rate limiter key with the Unix timestamp
	if meter == nil {
		meter = &iris_usage_meters.PayloadSizeMeter{}
	}
	rateLimiterKey := orgRateLimitTag(orgID)
	orgRequestsTotal := fmt.Sprintf("%d-total-zu-usage", orgID)
	if len(meter.Month) <= 0 {
		meter.Month = time.Now().UTC().Month().String()
	}
	orgRequestsMonthly := orgMonthlyUsageTag(orgID, meter.Month)
	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()

	// Increment the payload size meter
	_ = pipe.IncrByFloat(ctx, orgRequestsTotal, meter.ZeusRequestComputeUnitsConsumed())
	_ = pipe.IncrByFloat(ctx, orgRequestsMonthly, meter.ZeusRequestComputeUnitsConsumed())
	// Increment the rate limiter key
	_ = pipe.IncrByFloat(ctx, rateLimiterKey, meter.ZeusRequestComputeUnitsConsumed())

	// rate limiter key expires after 3 seconds
	pipe.Expire(ctx, rateLimiterKey, time.Second*3)

	// Generate the route key
	routeKey := orgRouteTag(orgID, rgName)

	// Check if the key exists
	existsCmd := pipe.Exists(ctx, routeKey)

	// Pop the endpoint from the head of the list and push it back to the tail, ensuring round-robin rotation
	endpointCmd := pipe.LMove(ctx, routeKey, routeKey, "LEFT", "RIGHT")

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
		return iris_models.RouteInfo{}, iris_usage_meters.UsageMeter{}, err
	} else if err == nil {
		// Convert the value to float
		rateLimitVal, err = strconv.ParseFloat(rateLimit, 64)
		if err != nil {
			fmt.Printf("error converting rate limit to float: %s\n", err.Error())
			return iris_models.RouteInfo{}, iris_usage_meters.UsageMeter{}, err
		}
	}

	monthlyUsage, err := monthlyUsageCmd.Result()
	monthlyUsageVal := 0.0
	if err != nil && err != redis.Nil {
		fmt.Printf("error getting monthly usage: %s\n", err.Error())
		return iris_models.RouteInfo{}, iris_usage_meters.UsageMeter{}, err
	} else if err == nil {
		// Convert the value to float
		monthlyUsageVal, err = strconv.ParseFloat(monthlyUsage, 64)
		if err != nil {
			fmt.Printf("error converting monthly usage to float: %s\n", err.Error())
			return iris_models.RouteInfo{}, iris_usage_meters.UsageMeter{}, err
		}
	}
	// Check whether the key exists
	if existsCmd.Val() <= 0 {
		log.Err(err).Msgf("key doesn't exist: %s", routeKey)
		//return iris_models.RouteInfo{}, fmt.Errorf("key doesn't exist: %s", routeKey)
	}
	// Get referers from Redis
	refererKey := orgRouteTag(orgID, endpointCmd.Val()) + ":referers"
	referers, err := m.Reader.SMembers(ctx, refererKey).Result()
	if err != nil {
		log.Err(err).Msgf("error getting referers for route: %s", refererKey)
		err = nil
		//return iris_models.RouteInfo{}, fmt.Errorf("error getting referers for route: %s", refererKey)
	}
	// Return the UsageRates
	return iris_models.RouteInfo{
			RoutePath: endpointCmd.Val(),
			Referers:  referers,
		}, iris_usage_meters.UsageMeter{
			RateLimit:    rateLimitVal,
			MonthlyUsage: monthlyUsageVal,
		}, nil
}

func (m *IrisCache) RecordRequestUsage(ctx context.Context, orgID int, meter *iris_usage_meters.PayloadSizeMeter) error {
	// Generate the rate limiter key with the Unix timestamp
	if meter == nil {
		meter = &iris_usage_meters.PayloadSizeMeter{}
	}
	rateLimiterKey := orgRateLimitTag(orgID)
	orgRequestsTotal := fmt.Sprintf("%d-total-zu-usage", orgID)
	if len(meter.Month) <= 0 {
		meter.Month = time.Now().UTC().Month().String()
	}
	orgRequestsMonthly := orgMonthlyUsageTag(orgID, meter.Month)
	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()
	// Increment the payload size meter
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
	return err
}

func (m *IrisCache) RecordRequestUsageRatesCheckLimitAndGetBroadcastRoutes(ctx context.Context, orgID int, procedureName, rgName string, meter *iris_usage_meters.PayloadSizeMeter) (iris_programmable_proxy_v1_beta.IrisRoutingProcedure, []iris_models.RouteInfo, iris_usage_meters.UsageMeter, error) {
	// Generate the rate limiter key with the Unix timestamp
	var procedureKey, procedureStepsKey string
	if orgID > 0 && procedureName != iris_programmable_proxy_v1_beta.MaxBlockAggReduce {
		procedureKey = fmt.Sprintf("%d:%s:procedure", orgID, procedureName)
		procedureStepsKey = fmt.Sprintf("%d:%s:procedure:steps", orgID, procedureName)
	} else {
		procedureKey = fmt.Sprintf("global:%s:procedure", procedureName)
		procedureStepsKey = fmt.Sprintf("global:%s:procedure:steps", procedureName)
	}

	rateLimiterKey := orgRateLimitTag(orgID)
	orgRequestsTotal := fmt.Sprintf("%d-total-zu-usage", orgID)
	if len(meter.Month) <= 0 {
		meter.Month = time.Now().UTC().Month().String()
	}
	orgRequestsMonthly := orgMonthlyUsageTag(orgID, meter.Month)
	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()

	// Increment the payload size meter
	_ = pipe.IncrByFloat(ctx, orgRequestsTotal, meter.ZeusRequestComputeUnitsConsumed())
	_ = pipe.IncrByFloat(ctx, orgRequestsMonthly, meter.ZeusRequestComputeUnitsConsumed())
	// Increment the rate limiter key
	_ = pipe.IncrByFloat(ctx, rateLimiterKey, meter.ZeusRequestComputeUnitsConsumed())

	// Get the values from Redis
	procedureCmd := pipe.Get(ctx, procedureKey)
	procedureStepsKeyCmd := pipe.Get(ctx, procedureStepsKey)
	if meter == nil {
		meter = &iris_usage_meters.PayloadSizeMeter{}
	}

	// rate limiter key expires after 3 seconds
	pipe.Expire(ctx, rateLimiterKey, time.Second*3)

	// Generate the route key
	routeKey := orgRouteTag(orgID, rgName)

	// Check if the key exists
	existsCmd := pipe.Exists(ctx, routeKey)

	endpointsCmd := pipe.LRange(ctx, routeKey, 0, -1)

	// Get the values from Redis
	rateLimitCmd := pipe.Get(ctx, rateLimiterKey)
	monthlyUsageCmd := pipe.Get(ctx, orgRequestsMonthly)

	// Execute the pipeline
	_, err := pipe.Exec(ctx)
	if err == redis.Nil {
		err = nil
	}
	var procedure iris_programmable_proxy_v1_beta.IrisRoutingProcedure
	// Get the values from the commands
	rateLimit, err := rateLimitCmd.Result()
	rateLimitVal := 0.0
	if err != nil && err != redis.Nil {
		fmt.Printf("error getting rate limit: %s\n", err.Error())
		return procedure, nil, iris_usage_meters.UsageMeter{}, err
	} else if err == nil {
		// Convert the value to float
		rateLimitVal, err = strconv.ParseFloat(rateLimit, 64)
		if err != nil {
			fmt.Printf("error converting rate limit to float: %s\n", err.Error())
			return procedure, nil, iris_usage_meters.UsageMeter{}, err
		}
	}
	monthlyUsage, err := monthlyUsageCmd.Result()
	monthlyUsageVal := 0.0
	if err != nil && err != redis.Nil {
		fmt.Printf("error getting monthly usage: %s\n", err.Error())
		return procedure, nil, iris_usage_meters.UsageMeter{}, err
	} else if err == nil {
		// Convert the value to float
		monthlyUsageVal, err = strconv.ParseFloat(monthlyUsage, 64)
		if err != nil {
			fmt.Printf("error converting monthly usage to float: %s\n", err.Error())
			return procedure, nil, iris_usage_meters.UsageMeter{}, err
		}
	}
	// Check whether the key exists
	if existsCmd.Val() <= 0 {
		log.Err(err).Msgf("key doesn't exist: %s", routeKey)
		//return iris_models.RouteInfo{}, fmt.Errorf("key doesn't exist: %s", routeKey)
	}
	// Get referers from Redis

	// Convert the endpoints to the desired type
	var routes []iris_models.RouteInfo
	for _, endpoint := range endpointsCmd.Val() {
		routeInfo := iris_models.RouteInfo{
			RoutePath: endpoint,
			Referers:  nil,
		}
		routes = append(routes, routeInfo)
	}

	data, err := procedureCmd.Bytes()
	if err != nil {
		log.Err(err).Msg("Failed to get procedure from Redis")
		return procedure, routes, iris_usage_meters.UsageMeter{}, err
	}
	// Deserialize the procedure
	err = json.Unmarshal(data, &procedure)
	if err != nil {
		log.Err(err).Msg("failed to deserialize the procedure")
		return procedure, routes, iris_usage_meters.UsageMeter{}, err
	}
	stepsBytes, err := procedureStepsKeyCmd.Bytes()
	if err != nil {
		log.Err(err).Msg("Failed to get procedure steps from Redis")
		return procedure, routes, iris_usage_meters.UsageMeter{}, err
	}
	var steps []iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep
	err = json.Unmarshal(stepsBytes, &steps)
	if err != nil {
		log.Err(err).Msg("Failed to deserialize procedure steps")
		return procedure, routes, iris_usage_meters.UsageMeter{}, err
	}
	for _, step := range steps {
		procedure.OrderedSteps.PushBack(step)
	}
	// Return the UsageRates
	return procedure, routes,
		iris_usage_meters.UsageMeter{
			RateLimit:    rateLimitVal,
			MonthlyUsage: monthlyUsageVal,
		}, nil
}