package hestia_web_router

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_delete "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/delete"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	hestia_access_keygen "github.com/zeus-fyi/olympus/hestia/web/access"
	hestia_billing "github.com/zeus-fyi/olympus/hestia/web/billing"
	hestia_quicknode_dashboard "github.com/zeus-fyi/olympus/hestia/web/iris"
	hestia_login "github.com/zeus-fyi/olympus/hestia/web/login"
	hestia_mev "github.com/zeus-fyi/olympus/hestia/web/mev"
	hestia_resources "github.com/zeus-fyi/olympus/hestia/web/resources"
	hestia_signup "github.com/zeus-fyi/olympus/hestia/web/signup"
	aegis_sessions "github.com/zeus-fyi/olympus/pkg/aegis/sessions"
)

func WebRoutes(e *echo.Echo) *echo.Echo {

	e.POST("/login", hestia_login.LoginHandler)
	e.POST("/discord/login", hestia_login.DiscordLoginHandler)
	e.GET("/reddit/callback", hestia_login.RedditLoginHandler)
	e.GET("/discord/callback", hestia_login.DiscordCallbackHandler)

	e.POST("/google/login", hestia_login.GoogleLoginHandler)
	e.POST("/signup", hestia_signup.SignUpHandler)
	e.GET("/logout/:token", Logout)
	e.GET("/v1/users/services", hestia_login.UsersServicesRequestHandler)

	e.GET("/twitter/callback", hestia_login.TwitterCallbackHandler)
	e.GET("/auth/:provider/callback", hestia_login.CallbackHandler)
	//e.GET("/logout/:provider", hestia_login.LogoutHandler)
	//e.GET("/auth/:provider", hestia_login.AuthHandler)

	e.GET("/verify/email/:token", hestia_signup.VerifyEmailHandler)
	hestia_quicknode_dashboard.InitQuickNodeDashboardRoutes(e)
	InitV1Routes(e)
	InitV1InternalRoutes(e)
	return e
}

const (
	QuickNodeMarketPlace        = "quickNodeMarketPlace"
	IrisQuickNodeTutorialToggle = "/quicknode/tutorial"
)

func InitV1Routes(e *echo.Echo) {
	eg := e.Group("/v1")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			cookie, err := c.Cookie(aegis_sessions.SessionIDNickname)
			if err == nil && cookie != nil {
				log.Info().Msg("InitV1ActionsRoutes: Cookie found")
				token = cookie.Value

			}
			key := read_keys.NewKeyReader()
			services, _, err := key.QueryUserAuthedServices(ctx, token)
			if err != nil {
				log.Err(err).Msg("InitV1Routes: QueryUserAuthedServices error")
				return false, err
			}
			c.Set("servicePlans", key.Services)
			ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
			c.Set("orgUser", ou)
			c.Set("bearer", key.PublicKey)
			return len(services) > 0, nil
		},
	}))
	eg.PUT(IrisQuickNodeTutorialToggle, hestia_quicknode_dashboard.TutorialToggleRequestHandler)
	eg.GET("/auth/status", hestia_access_keygen.AccessRequestHandler)
	eg.GET("/api/key/create", hestia_access_keygen.AccessKeyGenRequestHandler)
	eg.GET("/resources", hestia_resources.ResourceListRequestHandler)
	eg.GET("/stripe/customer/id", hestia_billing.StripeBillingRequestHandler)
	eg.GET("/refresh/token", hestia_login.TokenRefreshRequestHandler)

	eg.POST("/secrets/upsert", hestia_access_keygen.SecretsRequestHandler)
	eg.GET("/secrets", hestia_access_keygen.SecretsReadRequestHandler)
	eg.GET("/secret/:ref", hestia_access_keygen.SecretReadRequestHandler)
	eg.DELETE("/secret/:ref", hestia_access_keygen.SecretDeleteRequestHandler)
}

func InitV1InternalRoutes(e *echo.Echo) {
	eg := e.Group("/web/internal/v1")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			cookie, err := c.Cookie(aegis_sessions.SessionIDNickname)
			if err == nil && cookie != nil {
				log.Info().Msg("InitV1ActionsRoutes: Cookie found")
				token = cookie.Value
			}
			key, err := auth.VerifyInternalAdminBearerToken(ctx, token)
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
	eg.GET("/mev/dashboard", hestia_mev.MevRequestHandler)
}

func Logout(c echo.Context) error {
	ctx := context.Background()
	cookie, err := c.Cookie(aegis_sessions.SessionIDNickname)
	if err == nil && cookie != nil {
		log.Info().Msg("InitV1Routes: Cookie found")
		derr := hestia_delete.DeleteUserSessionKey(ctx, cookie.Value)
		if derr != nil {
			log.Err(derr).Msg("InitV1Routes: DeleteUserSessionKey error")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		return c.JSON(http.StatusOK, nil)
	}
	sessionToken := c.Param("token")
	derr := hestia_delete.DeleteUserSessionKey(ctx, sessionToken)
	if derr != nil {
		log.Err(derr).Msg("InitV1Routes: DeleteUserSessionKey error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}
