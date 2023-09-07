package v1internal_iris

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
)

const (
	RefreshAllOrgsRoutingTable = "/router/refresh/all"
	RefreshOrgRoutingTable     = "/router/refresh/:orgID"

	RefreshOrgGroupRoutingTable = "/router/refresh/:orgID/:groupName"
	DeleteOrgRoutingTables      = "/router/:orgID"
	DeleteOrgRoutingGroupTable  = "/router/:orgID/:groupName"

	DeleteQnOrgAuthCache = "/router/qn/auth/:qnID"

	DeleteSessionAuthCache = "/session/auth/:sessionID"
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

	eg.GET(RefreshAllOrgsRoutingTable, InternalRestoreCacheForAllOrgsHandler)
	eg.GET(RefreshOrgRoutingTable, InternalRefreshOrgRoutingTableHandler)
	eg.GET(RefreshOrgGroupRoutingTable, InternalRefreshOrgGroupRoutingTableHandler)

	eg.DELETE(DeleteOrgRoutingTables, InternalDeleteOrgRoutingTableRequestHandler)
	eg.DELETE(DeleteOrgRoutingGroupTable, InternalDeleteOrgGroupRoutingTableRequestHandler)
	eg.DELETE(DeleteQnOrgAuthCache, InternalDeleteQnOrgAuthCacheHandler)

	eg.DELETE(DeleteSessionAuthCache, InternalDeleteSessionAuthCacheHandler)
}
