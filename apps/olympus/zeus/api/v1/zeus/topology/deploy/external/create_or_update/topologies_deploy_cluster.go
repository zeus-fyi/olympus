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
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type TopologyClusterDeployRequest struct {
	ClusterClassName             string   `json:"clusterClassName"`
	SkeletonBaseOptions          []string `json:"skeletonBaseOptions"`
	AppTaint                     bool     `json:"appTaint,omitempty"`
	zeus_common_types.CloudCtxNs `json:"cloudCtxNs"`
}

func ClusterTopologyDeploymentHandler(c echo.Context) error {
	request := new(TopologyClusterDeployRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Interface("request", request).Msg("ClusterTopologyDeploymentHandler: orgUser")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, request.CloudCtxNs)
	if authed != true || err != nil {
		log.Err(err).Interface("ou", ou).Interface("request", request).Interface("authed", authed).Msg("ClusterTopologyDeploymentHandler: IsOrgCloudCtxNsAuthorized")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return request.DeployClusterTopology(c)
}

func (t *TopologyClusterDeployRequest) DeployClusterTopology(c echo.Context) error {
	log.Debug().Msg("DeployClusterTopology")
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	orgID := ou.OrgID
	cl, err := read_topology.SelectClusterTopology(ctx, orgID, t.ClusterClassName, t.SkeletonBaseOptions)
	if err != nil {
		log.Err(err).Msg("DeployClusterTopology: SelectClusterTopology")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	log.Info().Interface("cl", cl).Msg("DeployClusterTopology: SelectClusterTopology")
	clDeploy := base_deploy_params.ClusterTopologyWorkflowRequest{
		ClusterClassName:          t.ClusterClassName,
		TopologyIDs:               cl.GetTopologyIDs(),
		CloudCtxNS:                t.CloudCtxNs,
		OrgUser:                   ou,
		RequestChoreographySecret: cl.CheckForChoreographyOption(),
		AppTaint:                  t.AppTaint,
	}
	return zeus.ExecuteDeployClusterWorkflow(c, ctx, clDeploy)
}
