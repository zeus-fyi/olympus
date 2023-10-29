package iris_redis

import (
	"context"
	"time"

	"github.com/google/uuid"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

/*
AddOrUpdateOrgRoutingGroup

-> adds to routing tables using key:
rgTag := getOrgRouteKey(orgID, rgName)

data structures needed:

user -> sessions map
	getOrgSessionIDKey(orgID, sessionID) -> route

    1. needs add session + route
    2. needs get route via session id
	3. needs to remove anvil route while session is locked

	needs a ttl of some kind
	-> needs to delete session map to route when ttl expires
	-> needs to reset anvil + restore route to serverless router

rate limits
user -> sessions count
	getOrgSessionCountKey(orgID) -> count
	need to determine max active sessions per user

auto scaling up/down
	needs to trigger based on threshold low anvil servers in router

infra
	needs a class config

lb open questions
	- need to determine if a load balanced route can be used for anvil

needs user guide + examples

*/

func (r *IrisRedisTestSuite) TestSessions() {
	//rgName := "testGroupZ"
	routes := []iris_models.RouteInfo{
		{
			RoutePath: "https://zeus.fyi",
		},
		{
			RoutePath: "https://artemis.zeus.fyi",
		},
	}

	//anvilRoutes := []redis.Z{
	//	{
	//		Score:  1,
	//		Member: "https://anvil1.zeus.fyi",
	//	},
	//	{
	//		Score:  2,
	//		Member: "https://anvil2.zeus.fyi",
	//	},
	//}

	//err := IrisRedisClient.AddOrUpdateOrgRoutingGroup(context.Background(), 1, rgName, routes)
	//r.NoError(err)

	// AddOrUpdateOrgRoutingGroup
	orgID := 1
	for _, ep := range routes {
		rerr := IrisRedisClient.AddSessionWithTTL(context.Background(), orgID, uuid.New().String(), "http://serverless", ep.RoutePath, time.Second*10)
		r.NoError(rerr)
	}

	//m := make(map[string]int)
	//for i := 0; i < 10; i++ {
	//	routeEndpoint, rerr := IrisRedisClient.GetNextRoute(context.Background(), 1, rgName, nil)
	//	r.NoError(rerr)
	//	r.NotEmpty(routeEndpoint.RoutePath)
	//	fmt.Println(routeEndpoint)
	//	m[routeEndpoint.RoutePath]++
	//}
}
