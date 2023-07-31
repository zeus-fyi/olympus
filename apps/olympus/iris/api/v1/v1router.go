package v1_iris

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
)

var authCache = cache.New(time.Minute*5, time.Minute*10)

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
			val, ok := authCache.Get(token)
			if ok {
				orgUser, ok1 := val.(org_users.OrgUser)
				if ok1 {
					servicePlan, ok2 := authCache.Get(fmt.Sprintf("%d-servicePlan", orgUser.OrgID))
					if ok2 && len(servicePlan.(string)) > 0 {
						c.Set("servicePlan", servicePlan)
						c.Set("orgUser", orgUser)
						c.Set("bearer", val)
						return true, nil
					}
				}
			}
			key, err := auth.VerifyBearerTokenServiceWithQuickNodePlan(ctx, token, create_org_users.IrisQuickNodeService)
			if err != nil {
				log.Err(err).Msg("InitV1Routes")
				return false, c.JSON(http.StatusUnauthorized, nil)
			}
			if len(key.PublicKeyName) <= 0 {
				return false, c.JSON(http.StatusUnauthorized, Response{Message: "no service plan found"})
			}
			ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
			if key.PublicKeyVerified {
				authCache.SetDefault(fmt.Sprintf("%d-servicePlan", key.OrgID), key.PublicKeyName)
				authCache.SetDefault(token, key)
				authCache.SetDefault(fmt.Sprintf("%d", key.OrgID), ou)
			}
			c.Set("servicePlan", key.PublicKeyName)
			c.Set("orgUser", ou)
			c.Set("bearer", key.PublicKey)
			return key.PublicKeyVerified, err
		},
	}))
	eg.POST("/router/group", RpcLoadBalancerRequestHandler)
}
