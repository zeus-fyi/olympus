package destroy_deploy_request

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	topology_worker "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workers/topology"
	destroy_deployed_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/destroy"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/helpers"
	"go.temporal.io/sdk/client"
)

type TopologyDestroyDeployRequest struct {
	kns.TopologyKubeCtxNs
}

func (t *TopologyDestroyDeployRequest) DestroyDeployedTopology(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	tr := read_topology.NewInfraTopologyReaderWithOrgUser(ou)
	tr.TopologyID = t.TopologyID

	err := tr.SelectTopology(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// from auth lookup
	bearer := c.Get("bearer")
	tar := helpers.PackageCommonTopologyRequest(t.TopologyKubeCtxNs, bearer.(string), ou, tr.GetNativeK8s())

	workflowOptions := client.StartWorkflowOptions{}
	destroyDeployWf := destroy_deployed_workflow.NewDestroyDeployTopologyWorkflow()
	wf := destroyDeployWf.GetWorkflow()
	_, err = topology_worker.Worker.ExecuteWorkflow(ctx, workflowOptions, wf, tar)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("DestroyDeployedTopology, ExecuteWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := topology_deployment_status.NewTopologyStatus()

	resp.TopologyID = t.TopologyID
	resp.TopologyStatus = "Pending"
	resp.UpdatedAt = time.Now().UTC()
	return c.JSON(http.StatusAccepted, resp)
}
