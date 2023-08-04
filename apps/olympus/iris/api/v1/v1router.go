package v1_iris

import (
	"context"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
)

const QuickNodeMarketPlace = "quickNodeMarketPlace"

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
			orgID, plan, err := iris_redis.IrisRedisClient.GetAuthCacheIfExists(ctx, token)
			if err == nil && orgID > 0 && plan != "" {
				c.Set("servicePlan", plan)
				c.Set("orgUser", org_users.NewOrgUserWithID(int(orgID), 0))
				c.Set("bearer", token)
				return true, nil
			} else {
				plan = ""
				orgID = -1
				err = nil
			}
			key := read_keys.NewKeyReader()
			services, err := key.QueryUserAuthedServices(ctx, token)
			if err != nil {
				log.Err(err).Msg("InitV1Routes: QueryUserAuthedServices error")
				return false, err
			}
			c.Set("servicePlans", key.Services)
			if val, ok := key.Services[QuickNodeMarketPlace]; ok {
				plan = val
			} else {
				log.Warn().Str("marketplace", QuickNodeMarketPlace).Msg("InitV1Routes: marketplace not found")
				return false, errors.New("marketplace plan not found")
			}
			ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
			c.Set("servicePlan", plan)
			c.Set("orgUser", ou)
			c.Set("bearer", token)
			if err == nil && ou.OrgID > 0 && plan != "" {
				go func(oID int, token, plan string) {
					log.Info().Int("orgID", oID).Str("plan", plan).Msg("InitV1Routes: SetAuthCache")
					err = iris_redis.IrisRedisClient.SetAuthCache(context.Background(), oID, token, plan)
					if err != nil {
						log.Err(err).Msg("InitV1Routes: SetAuthCache")
					}
				}(ou.OrgID, token, plan)
			}
			return len(services) > 0, err
		},
	}))
	eg.POST("/router/group", RpcLoadBalancerPOSTRequestHandler)
	eg.GET("/router/group", RpcLoadBalancerGETRequestHandler)
	eg.PUT("/router/group", RpcLoadBalancerPUTRequestHandler)
	eg.DELETE("/router/group", RpcLoadBalancerDELETERequestHandler)
}
