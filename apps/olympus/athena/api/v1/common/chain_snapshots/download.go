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
	// download procedure
	pos := poseidon.NewPoseidon(athena.AthenaS3Manager)
	ctx := context.Background()
	pos.FnIn = t.ClientName + ".tar"
	pos.FnOut = t.ClientName
	err := pos.TarDownloadAndDec(ctx, t.BucketRequest)
	if err != nil {
		log.Err(err).Msg("DownloadChainSnapshotRequest: Lz4DownloadAndDec")
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = pos.RemoveFileInPath()
	if err != nil {
		log.Err(err).Msg("DownloadChainSnapshotRequest: Lz4DownloadAndDec")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
