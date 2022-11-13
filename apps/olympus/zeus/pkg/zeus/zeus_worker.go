package zeus

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	topology_worker "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workers/topology"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/helpers"
)

func ExecuteDeployWorkflow(c echo.Context, ctx context.Context, ou org_users.OrgUser, knsDeploy kns.TopologyKubeCtxNs, nk chart_workload.NativeK8s) error {
	tar := helpers.PackageCommonTopologyRequest(knsDeploy, ou, nk)
	err := topology_worker.Worker.ExecuteDeploy(ctx, tar)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("DeployTopology, ExecuteWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}

	resp := topology_deployment_status.NewTopologyStatus()
	resp.TopologyID = knsDeploy.TopologyID
	resp.TopologyStatus = topology_deployment_status.DeployPending
	resp.UpdatedAt = time.Now().UTC()
	return c.JSON(http.StatusAccepted, resp)
}
