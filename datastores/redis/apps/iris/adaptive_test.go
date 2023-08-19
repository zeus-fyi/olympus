package iris_redis

import (
	"context"
)

func (r *IrisRedisTestSuite) TestSetMetricLatencyTDigest() {
	err := IrisRedisClient.SetMetricLatencyTDigest(context.Background(), 1, "fooTable", "foo", 10.0)
	r.NoError(err)
}

func (r *IrisRedisTestSuite) TestGetMetricLatencyTDigest() {
	quantileVal, sc, err := IrisRedisClient.GetMetricPercentile(context.Background(), 1, "fooTable", "foo", 0.5)
	r.NoError(err)
	r.NotEmpty(quantileVal)
	r.NotEmpty(sc)
}
