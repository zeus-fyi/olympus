package create_infra

import (
	"fmt"

	"github.com/labstack/echo/v4"
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
	return nil
}
