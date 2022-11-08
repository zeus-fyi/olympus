package create_or_update_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	topology_worker "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workers/topology"
	deploy_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/create"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
	"go.temporal.io/sdk/client"
)

type TopologyDeployCreateActionDeployRequest struct {
	Action string
	base.TopologyActionRequest
}

func (t *TopologyDeployCreateActionDeployRequest) DeployTopology(c echo.Context) error {
	tr := read_topology.NewInfraTopologyReader()
	tr.TopologyID = t.TopologyID
	tr.OrgUser = t.OrgUser

	ctx := context.Background()
	err := tr.SelectTopology(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// TODO should also write kns & deployed status status
	kns := test.Kns
	kns.Namespace = "zeus"
	kns.CtxType = "dev-do-sfo3-zeus"
	workflowOptions := client.StartWorkflowOptions{}
	deployWf := deploy_workflow.NewDeployTopologyWorkflow()
	wf := deployWf.GetWorkflow()
	_, err = topology_worker.Worker.ExecuteWorkflow(ctx, workflowOptions, wf, t.TopologyActivityRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return err
}
