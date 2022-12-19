package create_or_update_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_infra "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/bases/infra"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

type TopologyClusterDeployRequest struct {
	ClusterClassName    string   `json:"clusterClassName"`
	SkeletonBaseOptions []string `json:"skeletonBaseOptions"`
	zeus_common_types.CloudCtxNs
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
	if t.ClusterClassName == "ethereumBeacons" {
		err := create_infra.InsertInfraBeaconCopy(ctx, ou)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("DeployClusterTopology: SelectClusterTopology")
		}
	}
	cl, err := read_topology.SelectClusterTopology(ctx, orgID, t.ClusterClassName, t.SkeletonBaseOptions)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("DeployClusterTopology: SelectClusterTopology")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	clDeploy := base_deploy_params.ClusterTopologyWorkflowRequest{
		ClusterName: t.ClusterClassName,
		TopologyIDs: cl.GetTopologyIDs(),
		CloudCtxNS:  t.CloudCtxNs,
		OrgUser:     ou,
	}
	return zeus.ExecuteDeployClusterWorkflow(c, ctx, clDeploy)
}
