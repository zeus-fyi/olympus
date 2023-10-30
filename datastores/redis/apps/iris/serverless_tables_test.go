package iris_redis

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

/*
infra
	needs a class config
lb open questions
	- need to determine if a load balanced route can be used for anvil

needs user guide + examples
*/

func (r *IrisRedisTestSuite) TestAddToServerlessTables() {
	anvilRoutes := []iris_models.RouteInfo{
		{
			RoutePath: "https://anvil1.zeus.fyi",
		},
		{
			RoutePath: "https://anvil2.zeus.fyi",
		},
	}

	err := IrisRedisClient.AddRoutesToServerlessRoutingTable(context.Background(), ServerlessAnvilTable, anvilRoutes)
	r.NoError(err)
}

func (r *IrisRedisTestSuite) TestGetServerlessTableRoutes() {
	sessionID := uuid.New().String()
	route, err := IrisRedisClient.GetNextServerlessRoute(context.Background(), 1, sessionID, ServerlessAnvilTable)
	r.NoError(err)
	r.NotEmpty(route)

	fmt.Println(route)
}

func (r *IrisRedisTestSuite) TestGetRouteFromSessionID() {
	sessionID := "94ce8571-411d-440b-a913-8a98d634894e"
	path, err := IrisRedisClient.GetServerlessSessionRoute(context.Background(), 1, sessionID)
	r.NoError(err)
	r.NotEmpty(path)
	fmt.Println(path)
}

func (r *IrisRedisTestSuite) TestServerlessRateLimit() {
	err := IrisRedisClient.CheckServerlessSessionRateLimit(context.Background(), 1, ServerlessAnvilTable)
	r.NoError(err)
}

func (r *IrisRedisTestSuite) TestReleaseServerlessRoute() {
	sessionID := "e6ed3599-c4fb-4ca5-8181-9cea50bcff21"
	err := IrisRedisClient.ReleaseServerlessRoute(context.Background(), 1, sessionID, ServerlessAnvilTable)
	r.NoError(err)
}
