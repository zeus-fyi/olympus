package iris_redis

import (
	"context"
	"fmt"

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

func (r *IrisRedisTestSuite) TestReadServerlessTableEntries() {
	routes, err := IrisRedisClient.GetServerlessRoutes(context.Background(), ServerlessAnvilTable)
	r.NoError(err)
	r.NotEmpty(routes)

	for _, route := range routes {
		fmt.Println(route)
	}
}

func (r *IrisRedisTestSuite) TestGetServerlessTableRoutes() {
	sessionID := "sessionID"
	route, _, err := IrisRedisClient.GetNextServerlessRoute(context.Background(), 1, sessionID, ServerlessAnvilTable)
	r.NoError(err)
	r.NotEmpty(route)

	fmt.Println(route)
}

func (r *IrisRedisTestSuite) TestGetRouteFromSessionID() {
	sessionID := "sessionID"
	path, err := IrisRedisClient.GetServerlessSessionRoute(context.Background(), 1, ServerlessAnvilTable, sessionID)
	r.NoError(err)
	r.NotEmpty(path)
	fmt.Println(path)
}

func (r *IrisRedisTestSuite) TestServerlessRateLimit() {
	_, err := IrisRedisClient.CheckServerlessSessionRateLimit(context.Background(), 1, "sessionID", ServerlessAnvilTable)
	r.NoError(err)
}

func (r *IrisRedisTestSuite) TestReleaseServerlessRoute() {
	sessionID := "sessionID"
	_, err := IrisRedisClient.ReleaseServerlessRoute(context.Background(), 1, sessionID, ServerlessAnvilTable)
	r.NoError(err)
}

func (r *IrisRedisTestSuite) TestPodNameExtract() {
	//pn, err := extractPodName("http://anvil-0.anvil.anvil-serverless-4d383226.svc.cluster.local:8545")
	//r.Require().NoError(err)
	//r.Assert().Equal("anvil-0", pn)
	//fmt.Println(pn)
}

func (r *IrisRedisTestSuite) TestRemoveAllServerlessTables() {
	sessionID := "sessionID"

	err := IrisRedisClient.DeleteAllServerlessTableAndOrgArtifacts(context.Background(), r.Tc.ProductionLocalTemporalOrgID, ServerlessAnvilTable, sessionID)
	r.NoError(err)
}
