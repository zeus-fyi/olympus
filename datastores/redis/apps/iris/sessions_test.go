package iris_redis

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func (r *IrisRedisTestSuite) TestSessions() {
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

	for _, ep := range routes {
		rerr := IrisRedisClient.AddSessionWithTTL(context.Background(), uuid.New().String(), ep.RoutePath, 1)
		r.NoError(rerr)
	}

	m := make(map[string]int)
	for i := 0; i < 10; i++ {
		routeEndpoint, rerr := IrisRedisClient.GetNextRoute(context.Background(), 1, rgName, nil)
		r.NoError(rerr)
		r.NotEmpty(routeEndpoint.RoutePath)
		fmt.Println(routeEndpoint)
		m[routeEndpoint.RoutePath]++
	}
}
