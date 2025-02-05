package zeus

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	topology_worker "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workers/topology"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/zeus/zeus/workload_config_drivers/topology_workloads"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func ExecuteCreateSetupClusterWorkflow(c echo.Context, ctx context.Context, params base_deploy_params.ClusterSetupRequest) error {
	err := topology_worker.Worker.ExecuteCreateSetupCluster(ctx, params)
	if err != nil {
		log.Err(err).Interface("orgUser", params.Ou.OrgID).Msg("ExecuteCreateSetupClusterWorkflow, ExecuteWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusAccepted, nil)
}

func ExecuteDeployFleetUpgradeWorkflow(c echo.Context, ctx context.Context, params base_deploy_params.FleetUpgradeWorkflowRequest) error {
	err := topology_worker.Worker.ExecuteDeployFleetUpgrade(ctx, params)
	if err != nil {
		log.Err(err).Interface("orgUser", params.OrgUser).Msg("ExecuteDeployFleetUpgradeWorkflow, ExecuteWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := topology_deployment_status.NewClusterTopologyStatus(params.ClusterName)
	resp.Status = topology_deployment_status.FleetRolloutRestartDeployPending
	return c.JSON(http.StatusAccepted, resp)
}

func ExecuteDeployFleetRolloutRestartWorkflow(c echo.Context, ctx context.Context, params base_deploy_params.FleetRolloutRestartWorkflowRequest) error {
	err := topology_worker.Worker.ExecuteDeployFleetRolloutRestart(ctx, params)
	if err != nil {
		log.Err(err).Interface("orgUser", params.OrgUser).Msg("ExecuteDeployFleetUpgradeWorkflow, ExecuteWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := topology_deployment_status.NewClusterTopologyStatus(params.ClusterName)
	resp.Status = topology_deployment_status.DeployPending
	return c.JSON(http.StatusAccepted, resp)
}

func ExecuteDeployRolloutRestartWorkflow(c echo.Context, ctx context.Context, ou org_users.OrgUser, cctxId int, cloudCtxNs zeus_common_types.CloudCtxNs) error {
	err := topology_worker.Worker.ExecuteDeployRolloutRestart(ctx, ou, cctxId, cloudCtxNs)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ExecuteDeployRolloutRestartWorkflow, ExecuteWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusAccepted, nil)
}

func ExecuteDeployClusterWorkflow(c echo.Context, ctx context.Context, params base_deploy_params.ClusterTopologyWorkflowRequest) error {
	err := topology_worker.Worker.ExecuteDeployCluster(ctx, params)
	if err != nil {
		log.Err(err).Interface("orgUser", params.OrgUser).Msg("ExecuteDeployClusterWorkflow, ExecuteWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := topology_deployment_status.NewClusterTopologyStatus(params.ClusterClassName)
	resp.Status = topology_deployment_status.DeployPending
	return c.JSON(http.StatusAccepted, resp)
}

func ExecuteDeployCronJobWorkflow(c echo.Context, ctx context.Context, ou org_users.OrgUser, tdr zeus_req_types.TopologyDeployRequest, nk topology_workloads.TopologyBaseInfraWorkload, deployChoreographySecret bool, clusterName, secretRef string) error {
	if nk.CronJob == nil && nk.ConfigMap == nil && nk.ServiceMonitor == nil {
		log.Err(nil).Interface("orgUser", ou).Interface("topologyID", tdr.TopologyID).Msg("ExecuteDeployCronJobWorkflow, payload is nil")
		return c.JSON(http.StatusBadRequest, nil)
	}
	tar := PackageCommonTopologyRequest(tdr, ou, nk, deployChoreographySecret, clusterName, secretRef)
	err := topology_worker.Worker.ExecuteDeployCronJob(ctx, tar)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("DeployCronJob, ExecuteWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := topology_deployment_status.NewTopologyStatus()
	resp.DeployStatus.TopologyID = tdr.TopologyID
	resp.TopologyStatus = topology_deployment_status.DeployPending
	resp.UpdatedAt = time.Now().UTC()
	return c.JSON(http.StatusAccepted, resp.DeployStatus)
}

func ExecuteDeployJobWorkflow(c echo.Context, ctx context.Context, ou org_users.OrgUser, tdr zeus_req_types.TopologyDeployRequest, nk topology_workloads.TopologyBaseInfraWorkload, deployChoreographySecret bool, clusterName, secretRef string) error {
	if nk.Job == nil && nk.ConfigMap == nil && nk.ServiceMonitor == nil {
		log.Err(nil).Interface("orgUser", ou).Interface("topologyID", tdr.TopologyID).Msg("ExecuteDeployJobWorkflow, payload is nil")
		return c.JSON(http.StatusBadRequest, nil)
	}
	tar := PackageCommonTopologyRequest(tdr, ou, nk, deployChoreographySecret, clusterName, secretRef)
	err := topology_worker.Worker.ExecuteDeployJob(ctx, tar)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("DeployJob, ExecuteWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := topology_deployment_status.NewTopologyStatus()
	resp.DeployStatus.TopologyID = tdr.TopologyID
	resp.TopologyStatus = topology_deployment_status.DeployPending
	resp.UpdatedAt = time.Now().UTC()
	return c.JSON(http.StatusAccepted, resp.DeployStatus)
}

func ExecuteDeployWorkflow(c echo.Context, ctx context.Context, ou org_users.OrgUser, knsDeploy zeus_req_types.TopologyDeployRequest, nk topology_workloads.TopologyBaseInfraWorkload, deployChoreographySecret bool, clusterName, secretRef string) error {
	if nk.Service == nil && nk.Deployment == nil && nk.StatefulSet == nil && nk.ServiceMonitor == nil && nk.Ingress == nil && nk.ConfigMap == nil {
		log.Err(nil).Interface("orgUser", ou).Interface("topologyID", knsDeploy.TopologyID).Msg("DeployTopology, payload is nil")
		return c.JSON(http.StatusBadRequest, nil)
	}
	tar := PackageCommonTopologyRequest(knsDeploy, ou, nk, deployChoreographySecret, clusterName, secretRef)
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

func ExecuteDestroyDeployWorkflow(c echo.Context, ctx context.Context, ou org_users.OrgUser, knsDestroyDeploy zeus_req_types.TopologyDeployRequest, nk topology_workloads.TopologyBaseInfraWorkload) error {
	tar := PackageCommonTopologyRequest(knsDestroyDeploy, ou, nk, false, "", "")
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

func ExecuteDestroyNamespaceWorkflow(c echo.Context, ctx context.Context, ou org_users.OrgUser, knsDestroyDeploy zeus_req_types.TopologyDeployRequest, nk topology_workloads.TopologyBaseInfraWorkload) error {
	tar := PackageCommonTopologyRequest(knsDestroyDeploy, ou, nk, false, "", "")
	err := topology_worker.Worker.ExecuteDestroyNamespace(ctx, tar)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ExecuteDestroyNamespace, ExecuteDestroyNamespaceWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := topology_deployment_status.NewTopologyStatus()
	resp.DeployStatus.TopologyID = knsDestroyDeploy.TopologyID
	resp.TopologyStatus = topology_deployment_status.DestroyDeployPending
	resp.UpdatedAt = time.Now().UTC()
	return c.JSON(http.StatusAccepted, resp.DeployStatus)
}

func ExecuteCleanDeployWorkflow(c echo.Context, ctx context.Context, ou org_users.OrgUser, knsCleanDeploy zeus_req_types.TopologyDeployRequest, nk topology_workloads.TopologyBaseInfraWorkload) error {
	tar := PackageCommonTopologyRequest(knsCleanDeploy, ou, nk, false, "", "")
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

func ExecuteDestroyResourcesWorkflow(c echo.Context, ctx context.Context, ou org_users.OrgUser, resources []int) error {
	dr := base_deploy_params.DestroyResourcesRequest{
		Ou:             ou,
		OrgResourceIDs: resources,
	}
	err := topology_worker.Worker.ExecuteDestroyResources(ctx, dr)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ExecuteDestroyResourcesWorkflow, ExecuteDestroyResources error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := topology_deployment_status.NewTopologyStatus()
	resp.TopologyStatus = topology_deployment_status.ResourceDestroyPending
	resp.UpdatedAt = time.Now().UTC()
	return c.JSON(http.StatusAccepted, resp.DeployStatus)
}
