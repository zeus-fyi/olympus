package iris_redis

import (
	"context"
	"fmt"

	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

/*
AddOrUpdateServerlessRoutingTable
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

	err := IrisRedisClient.AddOrUpdateServerlessRoutingTable(context.Background(), ServerlessAnvilTable, anvilRoutes)
	r.NoError(err)
}

func (r *IrisRedisTestSuite) TestGetServerlessTableRoutes() {
	route, err := IrisRedisClient.GetNextServerlessRoute(context.Background(), ServerlessAnvilTable)
	r.NoError(err)
	r.NotEmpty(route)

	fmt.Println(route)
}
