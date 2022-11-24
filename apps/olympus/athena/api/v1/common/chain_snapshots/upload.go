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

type UploadChainSnapshotRequest struct {
	poseidon.BucketRequest
}

func (t *UploadChainSnapshotRequest) Upload(c echo.Context) error {
	// upload procedure
	pos := poseidon.NewPoseidon(athena.AthenaS3Manager)
	ctx := context.Background()
	pos.DirIn = v1_common_routes.CommonManager.DataDir.DirIn
	err := pos.ZstdCompressAndUpload(ctx, t.BucketRequest)
	if err != nil {
		log.Err(err).Msg("UploadChainSnapshotRequest")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
