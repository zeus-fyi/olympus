package zeus

import (
	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_resp_types/topology_workloads"
)

func PackageCommonTopologyRequest(topCtxNs zeus_req_types.TopologyDeployRequest, ou org_users.OrgUser, nk topology_workloads.TopologyBaseInfraWorkload, deployChoreographySecret bool, clusterClassName, secretRef string) base_deploy_params.TopologyWorkflowRequest {
	topCtxNs.TopologyBaseInfraWorkload = nk
	topCtxNs.RequestChoreographySecretDeploy = deployChoreographySecret
	topCtxNs.ClusterClassName = clusterClassName
	topCtxNs.SecretRef = secretRef
	tar := base_deploy_params.TopologyWorkflowRequest{
		TopologyDeployRequest: topCtxNs,
		OrgUser:               ou,
	}
	return tar
}

func ConvertCloudCtxNsFormToType(c echo.Context) zeus_common_types.CloudCtxNs {
	cloudCtx := zeus_common_types.CloudCtxNs{}
	cloudCtx.CloudProvider = c.FormValue("cloudProvider")
	cloudCtx.Region = c.FormValue("region")
	cloudCtx.Context = c.FormValue("context")
	cloudCtx.Namespace = c.FormValue("namespace")
	cloudCtx.Env = c.FormValue("env")
	return cloudCtx
}
