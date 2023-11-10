package create_or_update_deploy

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
)

func TopologyDeploymentHandler(c echo.Context) error {
	request := new(TopologyDeployRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	// from auth lookup
	// validate context kns
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Msg("ClusterTopologyDeploymentHandler: orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, request.TopologyKubeCtxNs.CloudCtxNs)
	if authed != true || err != nil {
		log.Ctx(ctx).Err(err).Msg("ClusterTopologyDeploymentHandler: IsOrgCloudCtxNsAuthorized")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if request.TopologyKubeCtxNs.TopologyID == 0 {
		err = errors.New("no topology id provided")
		return c.JSON(http.StatusBadRequest, err)
	}
	return request.DeployTopology(c)
}
