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
	// upload procedure
	pos := poseidon.NewPoseidon(athena.AthenaS3Manager)
	ctx := context.Background()
	err := pos.ZstdCompressAndUpload(ctx, t.BucketRequest)
	if err != nil {
		log.Err(err).Msg("UploadChainSnapshotRequest")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
