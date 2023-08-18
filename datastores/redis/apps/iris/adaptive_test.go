package iris_redis

import "context"

func (r *IrisRedisTestSuite) TestAdaptiveCache() {
	err := IrisRedisClient.SetMetricLatencyTDigest(context.Background(), 1, "s", "test", 5.3)
	r.NoError(err)
}
