package v3_iris_internal

import (
	"context"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	v1_iris "github.com/zeus-fyi/olympus/iris/api/v1"
)

func InitV3InternalRoutes1(e *echo.Echo) {
	eg := e.Group("/v3/internal")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			orgU, plan, err := iris_redis.IrisRedisClient.GetInternalAuthCacheIfExists(ctx, token)
			if err == nil && orgU.OrgID > 0 && plan != "" {
				c.Set("lbDefault", v1_iris.GetDefaultLB(plan))
				c.Set("servicePlan", plan)
				c.Set("orgUser", org_users.NewOrgUserWithID(int(orgU.OrgID), orgU.UserID))
				return true, nil
			} else {
				return false, errors.New("internal auth cache not found")
			}
		},
	}))

	v1_iris.AddIrisRouter(eg)
}
