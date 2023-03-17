package zeus_router

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	aegis_sessions "github.com/zeus-fyi/olympus/pkg/aegis/sessions"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	zeus_v1_router "github.com/zeus-fyi/olympus/zeus/api/v1"
)

func InitRouter(e *echo.Echo, k8Cfg autok8s_core.K8Util) *echo.Echo {
	log.Debug().Msgf("InitRouter")
	// Routes
	e.GET("/health", Health)

	// external
	InitV1Routes(e, k8Cfg)
	// internal
	InitV1InternalRoutes(e, k8Cfg)
	// external
	InitV1ActionsRoutes(e, k8Cfg)
	return e
}

func InitV1ActionsRoutes(e *echo.Echo, k8Cfg autok8s_core.K8Util) {
	eg := e.Group("/v1")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			cookie, err := c.Cookie(aegis_sessions.SessionIDNickname)
			if err == nil && cookie != nil {
				token = cookie.Value
			}
			key, err := auth.VerifyBearerTokenService(ctx, token, create_org_users.ZeusService)
			if err != nil {
				log.Err(err).Msg("InitV1Routes")
				return false, c.JSON(http.StatusUnauthorized, nil)
			}
			ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
			c.Set("orgUser", ou)
			c.Set("bearer", key.PublicKey)
			return key.PublicKeyVerified, err
		},
	}))
	eg = zeus_v1_router.ActionsV1Routes(eg, k8Cfg)
}

func InitV1Routes(e *echo.Echo, k8Cfg autok8s_core.K8Util) {
	eg := e.Group("/v1")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			cookie, err := c.Cookie(aegis_sessions.SessionIDNickname)
			if err == nil && cookie != nil {
				token = cookie.Value
			}
			key, err := auth.VerifyBearerTokenService(ctx, token, create_org_users.ZeusService)
			if err != nil {
				log.Err(err).Msg("InitV1Routes")
				return false, c.JSON(http.StatusUnauthorized, nil)
			}
			ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
			c.Set("orgUser", ou)
			c.Set("bearer", key.PublicKey)
			return key.PublicKeyVerified, err
		},
	}))
	eg = zeus_v1_router.V1Routes(eg, k8Cfg)
}

func InitV1InternalRoutes(e *echo.Echo, k8Cfg autok8s_core.K8Util) {
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
	eg = zeus_v1_router.V1InternalRoutes(eg, k8Cfg)
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
