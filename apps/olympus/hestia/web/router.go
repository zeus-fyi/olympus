package hestia_web_router

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	hestia_delete "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/delete"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	hestia_login "github.com/zeus-fyi/olympus/hestia/web/login"
	hestia_signup "github.com/zeus-fyi/olympus/hestia/web/signup"
	aegis_sessions "github.com/zeus-fyi/olympus/pkg/aegis/sessions"
)

func WebRoutes(e *echo.Echo) *echo.Echo {
	e.POST("/login", hestia_login.LoginHandler)
	e.POST("/signup", hestia_signup.SignUpHandler)
	e.GET("/logout/:token", Logout)
	e.GET("/v1/users/services", hestia_login.UsersServicesRequestHandler)

	e.GET("/verify/email/:token", hestia_signup.VerifyEmailHandler)
	InitV1Routes(e)
	return e
}

func InitV1Routes(e *echo.Echo) {
	eg := e.Group("/v1")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
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
	eg.GET("/refresh/token", hestia_login.TokenRefreshRequestHandler)
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
