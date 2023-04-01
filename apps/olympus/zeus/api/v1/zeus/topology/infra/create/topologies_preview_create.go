package create_infra

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_resp_types/topology_workloads"
)

func PreviewCreateTopologyInfraActionRequestHandler(c echo.Context) error {
	request := new(TopologyPreviewCreateRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.PreviewCreateTopology(c)
}

type TopologyPreviewCreateRequest struct {
	Cluster `json:"cluster"`
}

func (t *TopologyPreviewCreateRequest) PreviewCreateTopology(c echo.Context) error {

	fmt.Println(t.Cluster)
	// TODO process
	tmp := topology_workloads.NewTopologyBaseInfraWorkload()
	return c.JSON(http.StatusOK, tmp)
}
