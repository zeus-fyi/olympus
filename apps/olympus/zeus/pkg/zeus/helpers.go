package zeus

import (
	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

func PackageCommonTopologyRequest(topCtxNs kns.TopologyKubeCtxNs, ou org_users.OrgUser, nk chart_workload.TopologyBaseInfraWorkload, deployChoreographySecret bool) base_deploy_params.TopologyWorkflowRequest {
	tar := base_deploy_params.TopologyWorkflowRequest{
		Kns:                       topCtxNs,
		OrgUser:                   ou,
		Host:                      "",
		TopologyBaseInfraWorkload: nk,
		RequestChoreographySecret: deployChoreographySecret,
	}
	return tar
}

func ConvertCloudCtxNsFormToType(c echo.Context) zeus_common_types.CloudCtxNs {
	cloudCtx := zeus_common_types.NewCloudCtxNs()
	cloudCtx.CloudProvider = c.FormValue("cloudProvider")
	cloudCtx.Region = c.FormValue("region")
	cloudCtx.Context = c.FormValue("context")
	cloudCtx.Namespace = c.FormValue("namespace")
	cloudCtx.Env = c.FormValue("env")
	return cloudCtx
}
