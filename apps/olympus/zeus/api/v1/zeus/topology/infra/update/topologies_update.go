package update_infra

import (
	"github.com/labstack/echo/v4"
	base_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/base"
)

type TopologyActionUpdateRequest struct {
	base_infra.TopologyInfraActionRequest
}

func (t *TopologyActionUpdateRequest) UpdateTopology(c echo.Context) error {
	return nil
}
