package external_api_workloads

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type TopologyCloudCtxNsQueryRequest struct {
	zeus_common_types.CloudCtxNs
}

func (t *TopologyCloudCtxNsQueryRequest) ReadDeployedWorkloads(c echo.Context) error {
	log.Debug().Msg("TopologyCloudCtxNsQueryRequest")
	ctx := context.Background()
	workload, err := zeus.K8Util.GetWorkloadAtNamespace(ctx, t.CloudCtxNs)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("TopologyCloudCtxNsQueryRequest: GetWorkloadAtNamespace")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, workload)
}
