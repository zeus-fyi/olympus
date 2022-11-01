package update_infra

import (
	"github.com/labstack/echo/v4"
)

type TopologyActionUpdateRequest struct {
}

func (t *TopologyActionUpdateRequest) UpdateTopology(c echo.Context) error {
	return nil
}
