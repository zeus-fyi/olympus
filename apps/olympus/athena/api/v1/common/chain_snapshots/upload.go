package athena_chain_snapshots

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"

	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
)

type UploadChainSnapshotRequest struct {
	poseidon.BucketRequest
}

type Response struct {
	Message string `json:"message"`
}

var ca = cache.New(15*time.Minute, 30*time.Minute)

func (t *UploadChainSnapshotRequest) Upload(c echo.Context) error {
	log.Info().Msg("UploadChainSnapshotRequest: Upload Sync Starting")
	pos := poseidon.NewPoseidon(athena.AthenaS3Manager)
	ctx := context.Background()
	resp := Response{Message: "done"}

	_, found := ca.Get(t.BucketRequest.GetBucketKey())
	if found {
		log.Info().Msg("UploadChainSnapshotRequest: Upload Sync Finished Recently Already")
		return c.JSON(http.StatusOK, resp)
	}
	err := pos.SyncUpload(ctx, t.BucketRequest)
	if err != nil {
		log.Err(err).Msg("Sync upload failed")
		return c.JSON(http.StatusInternalServerError, err)
	}
	log.Info().Msg("UploadChainSnapshotRequest: Upload Sync Finished")
	ca.Set(t.BucketRequest.GetBucketKey(), resp, 15*time.Minute)
	return c.JSON(http.StatusOK, resp)
}
