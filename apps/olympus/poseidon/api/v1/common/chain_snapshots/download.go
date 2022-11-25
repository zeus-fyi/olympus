package poseidon_chain_snapshots

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
)

type DownloadChainSnapshotRequest struct {
	poseidon.BucketRequest
}

func (t *DownloadChainSnapshotRequest) Download(c echo.Context) error {

	return c.JSON(http.StatusOK, nil)
}
