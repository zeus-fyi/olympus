package iris_redis

import "context"

func (r *IrisRedisTestSuite) TestSetMetricLatencyTDigest() {
	err := IrisRedisClient.SetMetricLatencyTDigest(context.Background(), 1, "s", "test", []float64{5.3}, false)
	r.NoError(err)
}

func (r *IrisRedisTestSuite) TestGetMetricPercentile() {
	// Assume that you have a method to mock the behavior of the Redis client
	// and it's assigned to IrisRedisClient.Reader
	value, err := IrisRedisClient.GetMetricPercentile(context.Background(), 1, "s", "test", 90.0)
	r.NoError(err)
	r.Equal(42.0, value)
}
