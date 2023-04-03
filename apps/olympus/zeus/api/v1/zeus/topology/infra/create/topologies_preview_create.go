package create_infra

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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
	//fmt.Println(t.Cluster)
	ctx := context.Background()
	pcg, err := zeus_templates.GenerateSkeletonBaseChartsPreview(ctx, t.Cluster)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error generating skeleton base charts")
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, pcg.ComponentBasesToSkeletonBases)
}
