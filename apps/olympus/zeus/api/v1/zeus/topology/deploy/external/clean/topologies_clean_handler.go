package clean_deploy_request

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func TopologyCleanNamespaceHandler(c echo.Context) error {
	request := new(zeus_req_types.TopologyDeployRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Interface("req", request).Msg("TopologyCleanNamespaceHandler, orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, request.CloudCtxNs)
	if authed != true {
		log.Warn().Interface("req", request).Interface("ou", ou).Msg("TopologyCleanNamespaceHandler, IsOrgCloudCtxNsAuthorized error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if err != nil {
		log.Warn().Interface("req", request).Interface("ou", ou).Msg("TopologyCleanNamespaceHandler, IsOrgCloudCtxNsAuthorized error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if request == nil {
		log.Warn().Interface("ou", ou).Msg("TopologyCleanNamespaceHandler, IsOrgCloudCtxNsAuthorized error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return CleanNamespaceAtTopology(c, *request)
}
