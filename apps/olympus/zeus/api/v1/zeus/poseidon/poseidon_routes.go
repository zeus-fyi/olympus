package snapshot_poseidon

import (
	"github.com/labstack/echo/v4"
)

func InternalSnapshotRoutes(e *echo.Group) *echo.Group {
	e.POST("/snapshot/sync", SnapshotHandler)
	return e
}
