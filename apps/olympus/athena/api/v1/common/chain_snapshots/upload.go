package athena_chain_snapshots

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
)

type UploadChainSnapshotRequest struct {
	poseidon.BucketRequest
}

func (t *UploadChainSnapshotRequest) Upload(c echo.Context) error {
	log.Info().Msg("UploadChainSnapshotRequest: Upload Sync Starting")
	pos := poseidon.NewPoseidon(athena.AthenaS3Manager)
	ctx := context.Background()
	err := pos.SyncUpload(ctx, t.BucketRequest)
	if err != nil {
		log.Err(err).Msg("Sync upload failed")
		return c.JSON(http.StatusInternalServerError, err)
	}
	log.Info().Msg("UploadChainSnapshotRequest: Upload Sync Finished")
	return c.JSON(http.StatusOK, nil)
}
