package iris_redis

import (
	"context"
	"fmt"
)

func (r *IrisRedisTestSuite) TestBroadcastRoutes() {
	rgName := "ethereum-mainnet"
	routes, err := IrisRedisClient.GetBroadcastRoutes(context.Background(), r.Tc.ProductionLocalTemporalOrgID, rgName)
	r.NoError(err)
	r.NotEmpty(routes)

	for _, ro := range routes {
		fmt.Println(ro.RoutePath)
	}
}
