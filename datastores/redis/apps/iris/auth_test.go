package iris_redis

import (
	"context"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (r *IrisRedisTestSuite) TestAuthCache() {
	ou := org_users.OrgUser{
		OrgUsers: autogen_bases.OrgUsers{UserID: 1, OrgID: 1},
	}
	err := IrisRedisClient.SetAuthCache(context.Background(), ou, "s", "enterprise", false)
	r.NoError(err)

	ou, plan, err := IrisRedisClient.GetAuthCacheIfExists(context.Background(), "test")
	r.NoError(err)
	r.Equal(int64(1), ou.OrgID)
	r.Equal("test", plan)
}
