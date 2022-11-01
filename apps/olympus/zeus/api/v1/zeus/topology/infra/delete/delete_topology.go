package delete_infra

import (
	"github.com/labstack/echo/v4"
)

type TopologyActionDeleteRequest struct {
}

func (t *TopologyActionDeleteRequest) DeleteTopology(c echo.Context) error {

	return nil
}
