package create_or_update_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type TopologyClusterDeployRequest struct {
	ClusterClassName             string   `json:"clusterClassName"`
	SkeletonBaseOptions          []string `json:"skeletonBaseOptions"`
	zeus_common_types.CloudCtxNs `json:"cloudCtxNs"`
}

func ClusterTopologyDeploymentHandler(c echo.Context) error {
	request := new(TopologyClusterDeployRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, request.CloudCtxNs)
	if authed != true || err != nil {
		log.Ctx(ctx).Err(err).Msg("ClusterTopologyDeploymentHandler: IsOrgCloudCtxNsAuthorized")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return request.DeployClusterTopology(c)
}

func (t *TopologyClusterDeployRequest) DeployClusterTopology(c echo.Context) error {
	log.Debug().Msg("DeployClusterTopology")
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	orgID := ou.OrgID
	cl, err := read_topology.SelectClusterTopology(ctx, orgID, t.ClusterClassName, t.SkeletonBaseOptions)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("DeployClusterTopology: SelectClusterTopology")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	log.Ctx(ctx).Info().Interface("cl", cl).Msg("DeployClusterTopology: SelectClusterTopology")
	clDeploy := base_deploy_params.ClusterTopologyWorkflowRequest{
		ClusterName:               t.ClusterClassName,
		TopologyIDs:               cl.GetTopologyIDs(),
		CloudCtxNS:                t.CloudCtxNs,
		OrgUser:                   ou,
		RequestChoreographySecret: cl.CheckForChoreographyOption(),
	}
	return zeus.ExecuteDeployClusterWorkflow(c, ctx, clDeploy)
}
