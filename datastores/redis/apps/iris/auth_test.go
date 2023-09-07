package iris_redis

import "context"

func (r *IrisRedisTestSuite) TestAuthCache() {
	err := IrisRedisClient.SetAuthCache(context.Background(), 1, "test", "test", false)
	r.NoError(err)

	orgID, plan, err := IrisRedisClient.GetAuthCacheIfExists(context.Background(), "test")
	r.NoError(err)
	r.Equal(int64(1), orgID)
	r.Equal("test", plan)
}
