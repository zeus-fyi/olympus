package iris_redis

import (
	"context"
	"fmt"

	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

func (r *IrisRedisTestSuite) TestDeleteOrgRoutingGroup() {
	rgName := "aptos-testnet"
	err := IrisRedisClient.DeleteOrgRoutingGroup(context.Background(), 1694205135884287000, rgName)
	r.NoError(err)
}

func (r *IrisRedisTestSuite) TestInitOrgTables() {
	rgName := "testGroupZ"
	routes := []iris_models.RouteInfo{
		{
			RoutePath: "https://zeus.fyi",
			Referrers: []string{"https://google.com", "https://yahoo.com"},
		},
		{
			RoutePath: "https://artemis.zeus.fyi",
			Referrers: nil,
		},
	}
	err := IrisRedisClient.AddOrUpdateOrgRoutingGroup(context.Background(), 1, rgName, routes)
	r.NoError(err)

	additionalRoute := "https://artemis2.zeus.fyi"
	routes = []iris_models.RouteInfo{
		{
			RoutePath: "https://zeus.fyi",
			Referrers: nil,
		},
		{
			RoutePath: "https://artemis.zeus.fyi",
			Referrers: nil,
		},
		{
			RoutePath: additionalRoute,
			Referrers: nil,
		},
	}
	err = IrisRedisClient.AddOrUpdateOrgRoutingGroup(context.Background(), 1, rgName, routes)
	r.NoError(err)
	m := make(map[string]int)
	for i := 0; i < 10; i++ {
		routeEndpoint, rerr := IrisRedisClient.GetNextRoute(context.Background(), 1, rgName, nil)
		r.NoError(rerr)
		r.NotEmpty(routeEndpoint)
		fmt.Println(routeEndpoint)
		m[routeEndpoint.RoutePath]++
	}

	routes = []iris_models.RouteInfo{
		{
			RoutePath: "https://zeus.fyi",
			Referrers: nil,
		},
		{
			RoutePath: "https://artemis.zeus.fyi",
			Referrers: nil,
		},
	}
	err = IrisRedisClient.AddOrUpdateOrgRoutingGroup(context.Background(), 1, rgName, routes)
	for i := 0; i < 10; i++ {
		routeEndpoint, rerr := IrisRedisClient.GetNextRoute(context.Background(), 1, rgName, nil)
		r.NoError(rerr)
		r.NotEmpty(routeEndpoint)
		fmt.Println(routeEndpoint)
		if routeEndpoint.RoutePath == additionalRoute {
			r.FailNow("route found in map")
		}
	}
}

func (r *IrisRedisTestSuite) TestRoundRobin() {
	rgName := "testGroupZ"
	routes := []iris_models.RouteInfo{
		{
			RoutePath: "https://zeus.fyi",
			Referrers: []string{"https://google.com", "https://yahoo.com"},
		},
		{
			RoutePath: "https://artemis.zeus.fyi",
			Referrers: nil,
		},
	}
	err := IrisRedisClient.AddOrUpdateOrgRoutingGroup(context.Background(), 1, rgName, routes)
	r.NoError(err)

	m := make(map[string]int)
	for i := 0; i < 10; i++ {
		routeEndpoint, rerr := IrisRedisClient.GetNextRoute(context.Background(), 1, rgName, nil)
		r.NoError(rerr)
		r.NotEmpty(routeEndpoint.RoutePath)
		fmt.Println(routeEndpoint)
		m[routeEndpoint.RoutePath]++
	}

	for _, ep := range routes {
		if m[ep.RoutePath] == 0 {
			r.FailNow("route not found in map")
		}
	}

	err = IrisRedisClient.DeleteOrgRoutingGroup(context.Background(), 1, rgName)
	r.NoError(err)

	_, rerr := IrisRedisClient.GetNextRoute(context.Background(), 1, rgName, nil)
	r.Error(rerr)
}

// DeleteOrgRoutingGroup
func (r *IrisRedisTestSuite) TestLoadBalancerRateMeter() {
	rgName := "testGroupZ"
	routes := []iris_models.RouteInfo{
		{
			RoutePath: "https://zeus.fyi",
			Referrers: []string{"https://google.com", "https://yahoo.com"},
		},
		{
			RoutePath: "https://artemis.zeus.fyi",
			Referrers: nil,
		},
	}
	err := IrisRedisClient.AddOrUpdateOrgRoutingGroup(context.Background(), 1, rgName, routes)
	r.NoError(err)

	meter := iris_usage_meters.NewPayloadSizeMeter(nil)

	m := make(map[string]int)
	for i := 0; i < 10; i++ {
		meter.Add(2048)
		routeEndpoint, rerr := IrisRedisClient.GetNextRoute(context.Background(), 1, rgName, meter)
		r.NoError(rerr)
		r.NotEmpty(routeEndpoint.RoutePath)
		fmt.Println(routeEndpoint)
		m[routeEndpoint.RoutePath]++

		meter.Reset()
		meter.Add(2048 * 10)
		err = IrisRedisClient.IncrementResponseUsageRateMeter(context.Background(), 1, meter)
		r.NoError(err)
		_, err = IrisRedisClient.CheckRateLimit(context.Background(), 1, "d", "test", meter)
		if err != nil {
			fmt.Println(err)
		}
	}

	meter.Reset()
	meter.Add(1024 * 1000)
	err = IrisRedisClient.IncrementResponseUsageRateMeter(context.Background(), 1, meter)
	r.NoError(err)
}

func (r *IrisRedisTestSuite) TestRateLimit() {
	meter := iris_usage_meters.NewPayloadSizeMeter(nil)

	ur, um, err := IrisRedisClient.RecordRequestUsageRatesCheckLimitAndNextRoute(context.Background(), 1, "d", meter)
	r.NoError(err)
	r.NotNil(um)
	r.NotEmpty(ur)
	fmt.Println(um)

	// should exceed monthly limit of 1k ZU
	_, err = IrisRedisClient.CheckRateLimit(context.Background(), 1, "test", "d", meter)
	r.Error(err)
}
