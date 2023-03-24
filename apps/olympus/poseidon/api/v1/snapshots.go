package v1_poseidon

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	pg_poseidon "github.com/zeus-fyi/olympus/datastores/postgres/apps/poseidon"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_orchestrations"
)

type SnapshotUploadRequest struct {
	pg_poseidon.UploadDataDirOrchestration
}

func SnapshotUploadRequestHandler(c echo.Context) error {
	log.Info().Msg("Poseidon: SnapshotUploadRequestHandler")
	request := new(SnapshotUploadRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("SnapshotUploadRequestHandler")
		return err
	}
	return request.ExecuteSnapshotUploadWorkflow(c)
}

func (s *SnapshotUploadRequest) ExecuteSnapshotUploadWorkflow(c echo.Context) error {
	ctx := context.Background()
	err := poseidon_orchestrations.PoseidonSyncWorker.ExecutePoseidonDiskUploadWorkflow(ctx, s.UploadDataDirOrchestration)
	if err != nil {
		log.Err(err).Msg("ExecuteSnapshotUploadWorkflow")
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusAccepted, nil)
}

type BeaconChainSyncUpload struct {
	ExecClient      pg_poseidon.UploadDataDirOrchestration
	ConsensusClient pg_poseidon.UploadDataDirOrchestration
}

func BeaconChainSyncUploadRequestHandler(c echo.Context) error {
	log.Info().Msg("Poseidon: BeaconChainSyncUploadRequestHandler")
	request := new(BeaconChainSyncUpload)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("SnapshotUploadRequestHandler")
		return err
	}
	return request.ExecuteBeaconChainSyncUploadWorkflow(c)
}

func (s *BeaconChainSyncUpload) ExecuteBeaconChainSyncUploadWorkflow(c echo.Context) error {
	ctx := context.Background()
	err := poseidon_orchestrations.PoseidonSyncWorker.ExecutePoseidonEthereumClientBeaconUploadWorkflow(ctx, s.ExecClient, s.ConsensusClient)
	if err != nil {
		log.Err(err).Msg("ExecuteBeaconChainSyncUploadWorkflow")
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusAccepted, nil)
}
