package poseidon_chain_snapshots

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
)

type UploadChainSnapshotRequest struct {
	poseidon.BucketRequest
}

func (t *UploadChainSnapshotRequest) Upload(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
