package v1_hera

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	ai_codegen "github.com/zeus-fyi/olympus/hera/api/v1/codegen"
	aegis_sessions "github.com/zeus-fyi/olympus/pkg/aegis/sessions"
)

func Routes(e *echo.Echo) *echo.Echo {
	e.GET("/health", Health)
	e.GET("/callback/oauth2", Health)
	InitV1BetaRoutes(e)
	return e
}

func InitV1BetaRoutes(e *echo.Echo) {
	eg := e.Group("/v1beta")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
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
			//c.Set("bearer", key.PublicKey)
			//b, err := hera_openai_dbmodels.CheckTokenBalance(ctx, ou)
			//if err != nil || b.TokensRemaining < 8000 {
			//	log.Err(err).Msg("InitV1BetaRoutes")
			//	if b.TokensRemaining < 8000 {
			//		return false, c.JSON(http.StatusBadRequest, "insufficient token balance, you need 8k min balance")
			//	}
			//	return false, c.JSON(http.StatusInternalServerError, "insufficient token balance, you need 8k min balance")
			//}
			return key.PublicKeyVerified, err
		},
	}))
	ai_codegen.CodeGenRoutes(eg)
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
