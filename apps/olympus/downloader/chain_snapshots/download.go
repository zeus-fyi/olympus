package hercules_chain_snapshots

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type DownloadChainSnapshotRequest struct {
	BucketRequest
}

func (t *DownloadChainSnapshotRequest) Download(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
