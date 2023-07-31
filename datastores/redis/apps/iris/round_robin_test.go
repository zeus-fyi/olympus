package iris_redis

import (
	"context"
	"fmt"

	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func (r *IrisRedisTestSuite) TestInitOrgTables() {
	rgName := "testGroupZ"
	routes := []iris_models.RouteInfo{
		{
			RoutePath: "https://zeus.fyi",
			Referers:  []string{"https://google.com", "https://yahoo.com"},
		},
		{
			RoutePath: "https://artemis.zeus.fyi",
			Referers:  nil,
		},
	}
	err := IrisRedis.AddOrUpdateOrgRoutingGroup(context.Background(), 1, rgName, routes)
	r.NoError(err)

	additionalRoute := "https://artemis2.zeus.fyi"
	routes = []iris_models.RouteInfo{
		{
			RoutePath: "https://zeus.fyi",
			Referers:  nil,
		},
		{
			RoutePath: "https://artemis.zeus.fyi",
			Referers:  nil,
		},
		{
			RoutePath: additionalRoute,
			Referers:  nil,
		},
	}
	err = IrisRedis.AddOrUpdateOrgRoutingGroup(context.Background(), 1, rgName, routes)
	r.NoError(err)
	m := make(map[string]int)
	for i := 0; i < 10; i++ {
		routeEndpoint, rerr := IrisRedis.GetNextRoute(context.Background(), 1, rgName)
		r.NoError(rerr)
		r.NotEmpty(routeEndpoint)
		fmt.Println(routeEndpoint)
		m[routeEndpoint.RoutePath]++
	}

	routes = []iris_models.RouteInfo{
		{
			RoutePath: "https://zeus.fyi",
			Referers:  nil,
		},
		{
			RoutePath: "https://artemis.zeus.fyi",
			Referers:  nil,
		},
	}
	err = IrisRedis.AddOrUpdateOrgRoutingGroup(context.Background(), 1, rgName, routes)
	for i := 0; i < 10; i++ {
		routeEndpoint, rerr := IrisRedis.GetNextRoute(context.Background(), 1, rgName)
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
			Referers:  []string{"https://google.com", "https://yahoo.com"},
		},
		{
			RoutePath: "https://artemis.zeus.fyi",
			Referers:  nil,
		},
	}
	err := IrisRedis.AddOrUpdateOrgRoutingGroup(context.Background(), 1, rgName, routes)
	r.NoError(err)

	m := make(map[string]int)
	for i := 0; i < 10; i++ {
		routeEndpoint, rerr := IrisRedis.GetNextRoute(context.Background(), 1, rgName)
		r.NoError(rerr)
		r.NotEmpty(routeEndpoint.RoutePath)
		fmt.Println(routeEndpoint)
		m[routeEndpoint.RoutePath]++
	}

	for _, ep := range routes {
		fmt.Printf("route: %s, count: %d\n", ep, m[ep.RoutePath])

		if m[ep.RoutePath] == 0 {
			r.FailNow("route not found in map")
		}
	}

	err = IrisRedis.DeleteOrgRoutingGroup(context.Background(), 1, rgName)
	r.NoError(err)

	_, rerr := IrisRedis.GetNextRoute(context.Background(), 1, rgName)
	r.Error(rerr)
}
