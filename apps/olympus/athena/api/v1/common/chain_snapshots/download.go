package athena_chain_snapshots

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	v1_common_routes "github.com/zeus-fyi/olympus/athena/api/v1/common"
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
	pos.DirOut = v1_common_routes.CommonManager.DataDir.DirIn
	err := pos.ZstdDownloadAndDec(ctx, t.BucketRequest)
	if err != nil {
		log.Err(err).Msg("DownloadChainSnapshotRequest")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
