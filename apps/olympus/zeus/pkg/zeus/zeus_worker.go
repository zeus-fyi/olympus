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
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
)

func ExecuteDeployClusterWorkflow(c echo.Context, ctx context.Context, params base_deploy_params.ClusterTopologyWorkflowRequest) error {
	err := topology_worker.Worker.ExecuteDeployCluster(ctx, params)
	if err != nil {
		log.Err(err).Interface("orgUser", params.OrgUser).Msg("ExecuteDeployClusterWorkflow, ExecuteWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := topology_deployment_status.NewClusterTopologyStatus(params.ClusterName)
	resp.Status = topology_deployment_status.DeployPending
	return c.JSON(http.StatusAccepted, resp)
}

func ExecuteDeployWorkflow(c echo.Context, ctx context.Context, ou org_users.OrgUser, knsDeploy kns.TopologyKubeCtxNs, nk chart_workload.TopologyBaseInfraWorkload) error {
	tar := PackageCommonTopologyRequest(knsDeploy, ou, nk)
	err := topology_worker.Worker.ExecuteDeploy(ctx, tar)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("DeployTopology, ExecuteWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := topology_deployment_status.NewTopologyStatus()
	resp.DeployStatus.TopologyID = knsDeploy.TopologyID
	resp.TopologyStatus = topology_deployment_status.DeployPending
	resp.UpdatedAt = time.Now().UTC()
	return c.JSON(http.StatusAccepted, resp.DeployStatus)
}

func ExecuteDestroyDeployWorkflow(c echo.Context, ctx context.Context, ou org_users.OrgUser, knsDestroyDeploy kns.TopologyKubeCtxNs, nk chart_workload.TopologyBaseInfraWorkload) error {
	tar := PackageCommonTopologyRequest(knsDestroyDeploy, ou, nk)
	err := topology_worker.Worker.ExecuteDestroyDeploy(ctx, tar)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("DestroyDeployedTopology, ExecuteWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := topology_deployment_status.NewTopologyStatus()
	resp.DeployStatus.TopologyID = knsDestroyDeploy.TopologyID
	resp.TopologyStatus = topology_deployment_status.DestroyDeployPending
	resp.UpdatedAt = time.Now().UTC()
	return c.JSON(http.StatusAccepted, resp.DeployStatus)
}

func ExecuteCleanDeployWorkflow(c echo.Context, ctx context.Context, ou org_users.OrgUser, knsCleanDeploy kns.TopologyKubeCtxNs, nk chart_workload.TopologyBaseInfraWorkload) error {
	tar := PackageCommonTopologyRequest(knsCleanDeploy, ou, nk)
	err := topology_worker.Worker.ExecuteCleanDeploy(ctx, tar)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ExecuteCleanDeploy, ExecuteWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := topology_deployment_status.NewTopologyStatus()
	resp.DeployStatus.TopologyID = knsCleanDeploy.TopologyID
	resp.TopologyStatus = topology_deployment_status.CleanDeployPending
	resp.UpdatedAt = time.Now().UTC()
	return c.JSON(http.StatusAccepted, resp.DeployStatus)
}
