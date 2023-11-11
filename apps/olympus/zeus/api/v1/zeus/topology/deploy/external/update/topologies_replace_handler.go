package deploy_updates

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func TopologyDeploymentReplaceHandler(c echo.Context) error {
	request := new(zeus_req_types.TopologyDeployRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	// from auth lookup
	// validate context kns
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Msg("TopologyDeploymentReplaceHandler, orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if request == nil {
		log.Warn().Interface("ou", ou).Msg("TopologyDeploymentReplaceHandler, request nil")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	topID := c.FormValue("topologyID")
	topologyID := string_utils.IntStringParser(topID)
	request.TopologyID = topologyID
	cloudCtxNs := zeus.ConvertCloudCtxNsFormToType(c)
	request.CloudCtxNs = cloudCtxNs
	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, cloudCtxNs)
	if authed != true {
		log.Warn().Interface("ou", ou).Msg("TopologyDeploymentReplaceHandler, request nil")
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err != nil {
		log.Err(err).Interface("ou", ou).Interface("req", request).Msg("TopologyDeploymentReplaceHandler, request nil")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return ReplaceTopology(c, *request)
}

func DeployUIClusterUpdateRequestHandler(c echo.Context) error {
	request := new(DeployClusterUpdateRequestUI)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.TopologyUpdateRequestUI(c)
}
