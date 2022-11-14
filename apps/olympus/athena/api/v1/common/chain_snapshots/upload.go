package athena_chain_snapshots

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UploadChainSnapshotRequest struct {
	ClientType string
	ClientName string
	Network    string
}

func (t *UploadChainSnapshotRequest) Upload(c echo.Context) error {
	// upload procedure
	return c.JSON(http.StatusOK, nil)
}
