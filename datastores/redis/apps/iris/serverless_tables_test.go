package iris_redis

import (
	"context"

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

	err := IrisRedisClient.AddOrUpdateServerlessRoutingTable(context.Background(), "anvil", anvilRoutes)
	r.NoError(err)
}
