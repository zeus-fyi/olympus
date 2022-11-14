package athena_chain_snapshots

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type DownloadChainSnapshotRequest struct {
	ClientType string
	ClientName string
	Network    string
}

func (t *DownloadChainSnapshotRequest) Download(c echo.Context) error {
	// download procedure
	return c.JSON(http.StatusOK, nil)
}
