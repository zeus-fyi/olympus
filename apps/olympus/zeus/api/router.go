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
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	zeus_v1_router "github.com/zeus-fyi/olympus/zeus/api/v1"
	zeus_webhooks "github.com/zeus-fyi/olympus/zeus/api/webhooks"
)

func InitRouter(e *echo.Echo, k8Cfg autok8s_core.K8Util, mw echo.MiddlewareFunc) *echo.Echo {
	log.Debug().Msgf("InitRouter")
	// Routes
	e.GET("/health", Health)

	// external
	InitV1Routes(e, k8Cfg, mw)
	// internal
	InitV1InternalRoutes(e, k8Cfg)

	InitV1WebhooksRoutes(e)
	InitVZWebhooksRoutes(e)
	// external
	InitV1ActionsRoutes(e, k8Cfg, mw)
	InitV1RoutesUI(e, k8Cfg, mw)
	return e
}

func InitV1ActionsRoutes(e *echo.Echo, k8Cfg autok8s_core.K8Util, mw echo.MiddlewareFunc) {
	eg := e.Group("/v1")
	eg.Use(mw, middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			cookie, err := c.Cookie(aegis_sessions.SessionIDNickname)
			if err == nil && cookie != nil {
				log.Info().Msg("InitV1ActionsRoutes: Cookie found")
				token = cookie.Value
			}
			key, err := auth.VerifyBearerTokenService(ctx, token, create_org_users.EthereumEphemeryService)
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

func InitV1Routes(e *echo.Echo, k8Cfg autok8s_core.K8Util, mw echo.MiddlewareFunc) {
	eg := e.Group("/v1")
	eg.Use(mw, middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			cookie, err := c.Cookie(aegis_sessions.SessionIDNickname)
			if err == nil && cookie != nil {
				log.Info().Msg("InitV1Routes: Cookie found")
				token = cookie.Value
			}
			key, err := auth.VerifyBearerTokenService(ctx, token, create_org_users.EthereumEphemeryService)
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

func InitV1RoutesUI(e *echo.Echo, k8Cfg autok8s_core.K8Util, mw echo.MiddlewareFunc) {
	eg := e.Group("/v1")
	eg.Use(mw, middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			cookie, err := c.Cookie(aegis_sessions.SessionIDNickname)
			if err == nil && cookie != nil {
				log.Info().Msg("InitV1Routes: Cookie found")
				token = cookie.Value
			}
			key, err := auth.VerifyBearerTokenService(ctx, token, create_org_users.EthereumEphemeryService)
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
	eg = zeus_v1_router.V1RoutesUI(eg, k8Cfg)
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

func InitV1WebhooksRoutes(e *echo.Echo) {
	eg := e.Group("/v1/webhooks")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			key, err := auth.VerifyInternalBearerToken(ctx, token)
			if err != nil {
				log.Err(err).Msg("InitV1InternalRoutes")
				return false, c.JSON(http.StatusInternalServerError, nil)
			}
			ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
			c.Set("orgUser", ou)
			c.Set("bearer", key.PublicKey)
			return key.PublicKeyVerified, err
		},
	}))
	eg.GET("/emails/support", zeus_webhooks.SupportEmailAIServiceTaskRequestHandler)
	eg.GET("/emails/alex", zeus_webhooks.AlexEmailAIServiceTaskRequestHandler)
	eg.GET("/emails/ai", zeus_webhooks.AiEmailAIServiceTaskRequestHandler)
	eg.GET("/telegram/ai/:group", zeus_webhooks.AiTelegramSupportAcknowledgeTelegramAiTaskHandler)
	eg.GET("/twitter/ai/:group", zeus_webhooks.SupportAcknowledgeTwitterAiTaskRequestHandler)
	eg.GET("/reddit/ai/:group", zeus_webhooks.SupportAcknowledgeRedditAiTaskRequestHandler)
	eg.GET("/discord/ai/:group", zeus_webhooks.SupportAcknowledgeDiscordAiTaskRequestHandler)
}

func InitVZWebhooksRoutes(e *echo.Echo) {
	eg := e.Group("/vz/webhooks")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			key, err := auth.FetchUserAuthTokenDiscord(ctx, 0)
			if err != nil {
				log.Err(err).Msg("InitV1InternalRoutes")
				return false, c.JSON(http.StatusInternalServerError, nil)
			}
			hs, err := misc.HashParams([]interface{}{key})
			if err != nil {
				log.Err(err).Msg("CreateDiscordJob: failed to hash params")
				return false, c.JSON(http.StatusInternalServerError, nil)
			}
			if token != hs {
				log.Err(err).Msg("InitV1InternalRoutes")
				return false, c.JSON(http.StatusInternalServerError, nil)
			}
			return true, err
		},
	}))
	eg.POST("/discord/ai", zeus_webhooks.RequestDiscordAiTaskStartRequestHandler)
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
