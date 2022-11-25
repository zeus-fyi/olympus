package poseidon_chain_snapshots

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	URL string
}

func (t *DownloadChainSnapshotRequest) GeneratePresignedURL(c echo.Context) error {
	resp := Response{}
	return c.JSON(http.StatusOK, resp)
}
