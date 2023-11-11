package zeus

import (
	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func PackageCommonTopologyRequest(topCtxNs zeus_req_types.TopologyDeployRequest, ou org_users.OrgUser, nk chart_workload.TopologyBaseInfraWorkload, deployChoreographySecret bool, clusterClassName, secretRef string) base_deploy_params.TopologyWorkflowRequest {
	tar := base_deploy_params.TopologyWorkflowRequest{
		TopologyDeployRequest:     topCtxNs,
		OrgUser:                   ou,
		RequestChoreographySecret: deployChoreographySecret,
		ClusterClassName:          clusterClassName,
		SecretRef:                 secretRef,
		TopologyBaseInfraWorkload: nk,
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
