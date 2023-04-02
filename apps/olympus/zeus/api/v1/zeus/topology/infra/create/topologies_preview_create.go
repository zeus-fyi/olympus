package create_infra

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
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
	zeus_templates.Cluster `json:"cluster"`
}

type TopologyPreviewCreateResponse struct {
	ComponentBases map[string]topology_workloads.TopologyBaseInfraWorkload `json:"componentBases"`
}

func (t *TopologyPreviewCreateRequest) PreviewCreateTopology(c echo.Context) error {

	fmt.Println(t.Cluster)
	/*
			// TODO process, map to ClusterPreview
			export interface ClusterPreview {
		    clusterName: string;
		    componentBases: any;
		    ingressSettings: any;
		    ingressPaths: any;
		}
	*/
	tmp := topology_workloads.NewTopologyBaseInfraWorkload()
	return c.JSON(http.StatusOK, tmp)
}
