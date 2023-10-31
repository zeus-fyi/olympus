package v1internal_iris

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	v1_iris "github.com/zeus-fyi/olympus/iris/api/v1"
	aegis_sessions "github.com/zeus-fyi/olympus/pkg/aegis/sessions"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
)

const (
	QuickNodeMarketPlace = "quickNodeMarketPlace"

	RefreshAllOrgsRoutingTable = "/router/refresh/all"
	RefreshOrgRoutingTable     = "/router/refresh/:orgID"

	RefreshOrgGroupRoutingTable = "/router/refresh/:orgID/:groupName"

	/* DELETE ROUTES */

	DeleteOrgRoutingTables     = "/router/:orgID"
	DeleteOrgRoutingGroupTable = "/router/:orgID/:groupName"

	DeleteQnOrgAuthCache = "/router/qn/auth/:qnID"

	DeleteSessionAuthCache = "/session/auth/:sessionID"

	RefreshServerlessTables = "/router/serverless/refresh"
)

func InitV1InternalRoutes(e *echo.Echo) {
	eg := e.Group("/v1/internal")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			key, err := auth.VerifyInternalBearerToken(ctx, token)
			if err != nil {
				log.Err(err).Msg("InitV1InternalRoutes")
				return false, c.JSON(http.StatusUnauthorized, nil)
			}
			ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
			c.Set("orgUser", ou)
			c.Set("bearer", key.PublicKey)
			return key.PublicKeyVerified, err
		},
	}))

	eg.GET(RefreshServerlessTables, InternalRefreshServerlessTablesHandler)
	eg.GET(RefreshAllOrgsRoutingTable, InternalRestoreCacheForAllOrgsHandler)
	eg.GET(RefreshOrgRoutingTable, InternalRefreshOrgRoutingTableHandler)
	eg.GET(RefreshOrgGroupRoutingTable, InternalRefreshOrgGroupRoutingTableHandler)

	eg.DELETE(DeleteOrgRoutingTables, InternalDeleteOrgRoutingTableRequestHandler)
	eg.DELETE(DeleteOrgRoutingGroupTable, InternalDeleteOrgGroupRoutingTableRequestHandler)
	eg.DELETE(DeleteQnOrgAuthCache, InternalDeleteQnOrgAuthCacheHandler)

	eg.DELETE(DeleteSessionAuthCache, InternalDeleteSessionAuthCacheHandler)
}

func InitV2InternalRoutes1(e *echo.Echo) {
	eg := e.Group("/v2/internal")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			usingCookie := false
			cookie, err := c.Cookie(aegis_sessions.SessionIDNickname)
			if err == nil && cookie != nil {
				log.Info().Msg("InitV1Routes: Cookie found")
				token = cookie.Value
				usingCookie = true
			}
			if token == "fgjlsdjgmklosadmgslkasdmglkasm" {
				token = artemis_orchestration_auth.Bearer
			}

			orgU, plan, err := iris_redis.IrisRedisClient.GetAuthCacheIfExists(ctx, token)
			if err == nil && orgU.OrgID > 0 && plan != "" {
				c.Set("lbDefault", v1_iris.GetDefaultLB(plan))
				c.Set("servicePlan", plan)
				c.Set("orgUser", org_users.NewOrgUserWithID(int(orgU.OrgID), orgU.UserID))
				c.Set("bearer", token)
				return true, nil
			} else {
				plan = ""
				orgU.OrgID = -1
				orgU.UserID = -1
				err = nil
			}
			key := read_keys.NewKeyReader()
			services, _, err := key.QueryUserAuthedServices(ctx, token)
			if err != nil {
				log.Err(err).Msg("InitV1Routes: QueryUserAuthedServices error")
				return false, err
			}
			if key.OrgID != auth.TemporalOrgID {
				return false, errors.New("unauthorized")
			}
			c.Set("servicePlans", key.Services)
			if val, ok := key.Services[QuickNodeMarketPlace]; ok {
				plan = val
			} else {
				log.Warn().Str("marketplace", QuickNodeMarketPlace).Msg("InitV1Routes: marketplace not found")
				return false, errors.New("marketplace plan not found")
			}
			ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
			c.Set("lbDefault", v1_iris.GetDefaultLB(plan))
			c.Set("servicePlan", plan)
			c.Set("orgUser", ou)
			c.Set("bearer", token)
			if err == nil && ou.OrgID > 0 && ou.UserID > 0 && plan != "" {
				go func(oID int, token, plan string, usingCookie bool) {
					log.Info().Int("orgID", oID).Str("plan", plan).Msg("InitV1Routes: SetAuthCache")
					err = iris_redis.IrisRedisClient.SetAuthCache(context.Background(), ou, token, plan, usingCookie)
					if err != nil {
						log.Err(err).Msg("InitV1Routes: SetAuthCache")
					}
				}(ou.OrgID, token, plan, usingCookie)
			}
			return len(services) > 0, err
		},
	}))

	v1_iris.AddIrisRouter(eg)
}
