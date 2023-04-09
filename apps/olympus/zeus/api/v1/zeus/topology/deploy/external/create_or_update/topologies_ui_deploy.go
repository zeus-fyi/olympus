package create_or_update_deploy

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

//const payload = {
//"node": node,
//"count": count,
//"namespaceAlias": namespaceAlias,
//"cluster": cluster,
//"resourceRequirements": resourceRequirements,
//}

type TopologyDeployUIRequest struct {
	zeus_common_types.CloudCtxNs
	Node                 any `json:"node"`
	Count                any `json:"count"`
	NamespaceAlias       any `json:"namespaceAlias"`
	Cluster              any `json:"cluster"`
	ResourceRequirements any `json:"resourceRequirements"`
}

func SetupClusterTopologyDeploymentHandler(c echo.Context) error {
	request := new(TopologyDeployUIRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.DeploySetupClusterTopology(c)
}

func (t *TopologyDeployUIRequest) DeploySetupClusterTopology(c echo.Context) error {
	log.Debug().Msg("DeploySetupClusterTopology")
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	clusterID := uuid.New()
	cr := base_deploy_params.ClusterSetupRequest{
		Ou:            ou,
		CloudCtxNs:    zeus_common_types.CloudCtxNs{},
		ClusterID:     clusterID,
		Nodes:         autogen_bases.Nodes{},
		NodesQuantity: t.Count.(float64),
		Disks:         autogen_bases.Disks{},
		DisksQuantity: 0, // todo
	}
	return zeus.ExecuteCreateSetupClusterWorkflow(c, ctx, cr)
}
