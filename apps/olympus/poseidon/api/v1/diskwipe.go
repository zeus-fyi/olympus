package v1_poseidon

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	pg_poseidon "github.com/zeus-fyi/olympus/datastores/postgres/apps/poseidon"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_orchestrations"
)

type DiskWipeRequest struct {
	pg_poseidon.DiskWipeOrchestration
}

func DiskWipeRequestHandler(c echo.Context) error {
	log.Info().Msg("Poseidon: DiskWipeRequestHandler")
	request := new(DiskWipeRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("DiskWipeRequestHandler")
		return err
	}
	return request.ExecuteDiskWipeWorkflow(c)
}

func (dw *DiskWipeRequest) ExecuteDiskWipeWorkflow(c echo.Context) error {
	ctx := context.Background()
	err := poseidon_orchestrations.PoseidonSyncWorker.ExecutePoseidonDiskWipeWorkflow(ctx, dw.DiskWipeOrchestration)
	if err != nil {
		log.Err(err).Msg("ExecuteDiskWipeWorkflow")
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusAccepted, nil)
}
