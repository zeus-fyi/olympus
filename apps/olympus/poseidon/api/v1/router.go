package v1_poseidon

import (
	"net/http"

	"github.com/labstack/echo/v4"
	poseidon_chain_snapshots "github.com/zeus-fyi/olympus/poseidon/api/v1/common/chain_snapshots"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Routes
	e.GET("/health", Health)

	e.POST("/snapshot/download", poseidon_chain_snapshots.RequestDownloadChainSnapshotHandler)
	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
