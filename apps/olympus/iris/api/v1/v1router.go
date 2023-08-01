package v1_iris

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
)

func InitV1Routes(e *echo.Echo) {
	eg := e.Group("/v1")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			// Get headers
			qnTestHeader := c.Request().Header.Get(QuickNodeTestHeader)
			qnIDHeader := c.Request().Header.Get(QuickNodeIDHeader)
			qnEndpointID := c.Request().Header.Get(QuickNodeEndpointID)
			qnChain := c.Request().Header.Get(QuickNodeChain)
			qnNetwork := c.Request().Header.Get(QuickNodeNetwork)

			// Set headers to echo context
			c.Set(QuickNodeTestHeader, qnTestHeader)
			c.Set(QuickNodeIDHeader, qnIDHeader)
			c.Set(QuickNodeEndpointID, qnEndpointID)
			c.Set(QuickNodeChain, qnChain)
			c.Set(QuickNodeNetwork, qnNetwork)
			orgID, plan, err := iris_redis.IrisRedis.GetAuthCacheIfExists(ctx, token)
			if err != nil && orgID > 0 && plan != "" {
				c.Set("servicePlan", plan)
				c.Set("orgUser", org_users.NewOrgUserWithID(int(orgID), 0))
				c.Set("bearer", token)
				return true, nil
			} else {
				plan = ""
				orgID = -1
				err = nil
			}
			key, err := auth.VerifyBearerTokenServiceWithQuickNodePlan(ctx, token, create_org_users.IrisQuickNodeService)
			if err != nil {
				log.Err(err).Msg("InitV1Routes")
				return false, c.JSON(http.StatusUnauthorized, nil)
			}
			if len(key.PublicKeyName) <= 0 {
				return false, c.JSON(http.StatusUnauthorized, Response{Message: "no service plan found"})
			}
			if err == nil && key.PublicKeyVerified && orgID > 0 && plan != "" {
				go func(oID int, token, plan string) {
					err = iris_redis.IrisRedis.SetAuthCache(context.Background(), oID, token, key.PublicKeyName)
					if err != nil {
						log.Err(err).Msg("InitV1Routes: SetAuthCache")
					}
				}(int(orgID), token, key.PublicKeyName)
				ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
				c.Set("servicePlan", plan)
				c.Set("orgUser", ou)
				c.Set("bearer", token)
			}
			return key.PublicKeyVerified, err
		},
	}))
	eg.POST("/router/group", RpcLoadBalancerRequestHandler)
}
