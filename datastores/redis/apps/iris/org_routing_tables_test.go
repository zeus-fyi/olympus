package iris_redis

import (
	"context"
)

func (r *IrisRedisTestSuite) TestInitRoutingTables() {

	err := IrisRedisClient.initRoutingTables(context.Background())
	r.Require().Nil(err)
}
