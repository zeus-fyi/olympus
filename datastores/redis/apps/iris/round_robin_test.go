package iris_redis

import (
	"context"
	"fmt"
)

func (r *IrisRedisTestSuite) TestInitOrgTables() {
	rgName := "testGroupZ"
	routes := []string{"https://zeus.fyi", "https://artemis.zeus.fyi"}
	err := IrisRedis.AddOrUpdateOrgRoutingGroup(context.Background(), 1, rgName, routes)
	r.NoError(err)

	additionalRoute := "https://artemis2.zeus.fyi"
	routes = []string{"https://zeus.fyi", "https://artemis.zeus.fyi", additionalRoute}

	err = IrisRedis.AddOrUpdateOrgRoutingGroup(context.Background(), 1, rgName, routes)
	r.NoError(err)
	m := make(map[string]int)
	for i := 0; i < 10; i++ {
		routeEndpoint, rerr := IrisRedis.GetNextRoute(context.Background(), 1, rgName)
		r.NoError(rerr)
		r.NotEmpty(routeEndpoint)
		fmt.Println(routeEndpoint)
		m[routeEndpoint]++
	}

	routes = []string{"https://zeus.fyi", "https://artemis.zeus.fyi"}
	err = IrisRedis.AddOrUpdateOrgRoutingGroup(context.Background(), 1, rgName, routes)
	for i := 0; i < 10; i++ {
		routeEndpoint, rerr := IrisRedis.GetNextRoute(context.Background(), 1, rgName)
		r.NoError(rerr)
		r.NotEmpty(routeEndpoint)
		fmt.Println(routeEndpoint)
		if routeEndpoint == additionalRoute {
			r.FailNow("route found in map")
		}
	}
}

func (r *IrisRedisTestSuite) TestRoundRobin() {
	rgName := "testGroupZ"
	routes := []string{"https://zeus.fyi", "https://artemis.zeus.fyi"}

	m := make(map[string]int)
	for i := 0; i < 10; i++ {
		routeEndpoint, rerr := IrisRedis.GetNextRoute(context.Background(), 1, rgName)
		r.NoError(rerr)
		r.NotEmpty(routeEndpoint)
		fmt.Println(routeEndpoint)
		m[routeEndpoint]++
	}

	for _, ep := range routes {
		fmt.Printf("route: %s, count: %d\n", ep, m[ep])

		if m[ep] == 0 {
			r.FailNow("route not found in map")
		}
	}

	err := IrisRedis.DeleteOrgRoutingGroup(context.Background(), 1, rgName)
	r.NoError(err)

	_, rerr := IrisRedis.GetNextRoute(context.Background(), 1, rgName)
	r.Error(rerr)
}
