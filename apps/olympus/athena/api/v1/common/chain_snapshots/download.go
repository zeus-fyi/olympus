package athena_chain_snapshots

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
)

type DownloadChainSnapshotRequest struct {
	poseidon.BucketRequest
}

func (t *DownloadChainSnapshotRequest) Download(c echo.Context) error {
	log.Info().Msg("DownloadChainSnapshotRequest: Download Sync Starting")
	pos := poseidon.NewPoseidon(athena.AthenaS3Manager)
	ctx := context.Background()
	err := pos.SyncDownload(ctx, t.BucketRequest)
	if err != nil {
		log.Err(err).Msg("Sync")
		return c.JSON(http.StatusInternalServerError, err)
	}
	log.Info().Msg("DownloadChainSnapshotRequest: Download Sync Finished")
	return c.JSON(http.StatusOK, nil)
}
