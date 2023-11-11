package create_or_update_deploy

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func TopologyDeploymentHandler(c echo.Context) error {
	request := new(zeus_req_types.TopologyDeployRequest)
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
	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, request.CloudCtxNs)
	if authed != true || err != nil {
		if err == nil {
			err = fmt.Errorf("not authorized")
		}
		log.Err(err).Interface("ou", ou).Interface("req", request).Msg("ClusterTopologyDeploymentHandler: IsOrgCloudCtxNsAuthorized")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if request.TopologyID == 0 {
		err = errors.New("no topology id provided")
		log.Err(err).Interface("ou", ou).Interface("req", request).Msg("ClusterTopologyDeploymentHandler: IsOrgCloudCtxNsAuthorized")
		return c.JSON(http.StatusBadRequest, err)
	}
	if request == nil {
		log.Warn().Interface("ou", ou).Interface("req", request).Msg("ClusterTopologyDeploymentHandler: IsOrgCloudCtxNsAuthorized")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return DeployTopology(c, *request)
}
