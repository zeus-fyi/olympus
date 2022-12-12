package v1_poseidon

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_orchestrations"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Routes
	e.GET("/health", Health)
	InitV1Routes(e)
	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}

func InitV1Routes(e *echo.Echo) {
	eg := e.Group("/v1beta")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			key, err := auth.VerifyBearerToken(ctx, token)
			if err != nil {
				log.Err(err).Msg("InitV1Routes")
				return false, c.JSON(http.StatusInternalServerError, nil)
			}
			ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
			c.Set("orgUser", ou)
			c.Set("bearer", key.PublicKey)
			return key.PublicKeyVerified, err
		},
	}))
	eg.GET("/ethereum", RunWorkflow)
}

func RunWorkflow(c echo.Context) error {
	return ExecuteSyncWorkflow(c, context.Background())
}

func ExecuteSyncWorkflow(c echo.Context, ctx context.Context) error {
	err := poseidon_orchestrations.PoseidonSyncWorker.ExecutePoseidonSyncWorkflow(ctx)
	if err != nil {
		log.Err(err).Msg("ExecuteSyncWorkflow")
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusAccepted, nil)
}
