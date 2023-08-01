package iris_redis

import "context"

func (r *IrisRedisTestSuite) TestAuthCache() {
	err := IrisRedis.SetAuthCache(context.Background(), 1, "test", "test")
	r.NoError(err)

	orgID, plan, err := IrisRedis.GetAuthCacheIfExists(context.Background(), "test")
	r.NoError(err)
	r.Equal(int64(1), orgID)
	r.Equal("test", plan)
}
