package v1_athena

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	v1_common_routes "github.com/zeus-fyi/olympus/athena/api/v1/common"
	athena_chain_snapshots "github.com/zeus-fyi/olympus/athena/api/v1/common/chain_snapshots"
	host "github.com/zeus-fyi/olympus/athena/api/v1/common/host_info"
	athena_jwt_route "github.com/zeus-fyi/olympus/athena/api/v1/common/jwt"
	athena_routines "github.com/zeus-fyi/olympus/athena/api/v1/common/routines"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func InitV1InternalRoutes(e *echo.Echo, p filepaths.Path) {
	eg := e.Group("/v1/internal")
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
	eg = CommonRoutes(eg, p)
	return
}

func CommonRoutes(e *echo.Group, p filepaths.Path) *echo.Group {
	v1_common_routes.CommonManager.DataDir = p
	e.POST("/jwt/create", athena_jwt_route.JwtHandler)

	e.POST("/routines/suspend", athena_routines.SuspendRoutineHandler)
	e.POST("/routines/start", athena_routines.StartAppRoutineHandler)
	e.POST("/routines/resume", athena_routines.ResumeProcessRoutineHandler)
	e.POST("/routines/kill", athena_routines.KillProcessRoutineHandler)

	e.POST("/routines/disk/wipe", athena_routines.WipeDiskHandler)

	e.GET("/host/disk", host.GetDiskStatsHandler)
	e.GET("/host/memory", host.GetMemStatsHandler)

	e.POST("/snapshot/download", athena_chain_snapshots.DownloadChainSnapshotHandler)
	e.POST("/snapshot/upload", athena_chain_snapshots.UploadChainSnapshotHandler)

	return e
}
