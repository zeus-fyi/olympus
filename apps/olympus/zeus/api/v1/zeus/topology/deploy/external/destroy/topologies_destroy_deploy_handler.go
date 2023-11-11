package destroy_deploy_request

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func TopologyDestroyDeploymentHandler(c echo.Context) error {
	request := new(zeus_req_types.TopologyDeployRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Msg("TopologyDestroyDeploymentHandler: orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if request == nil {
		log.Warn().Interface("ou", ou).Msg("TopologyDestroyDeploymentHandler: request empty error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, request.CloudCtxNs)
	if authed != true {
		log.Warn().Interface("ou", ou).Interface("req", request).Msg("TopologyDestroyDeploymentHandler: IsOrgCloudCtxNsAuthorized error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("TopologyDestroyDeploymentHandler: IsOrgCloudCtxNsAuthorized error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return DestroyDeployedTopology(c, *request)
}

func DestroyResourceHandler(c echo.Context) error {
	request := new(ResourceDestroyRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.DestroyResource(c)
}

func DestroyNamespaceHandler(c echo.Context) error {
	request := new(TopologyUIDestroyDeployRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("DestroyNamespaceHandler: Bind error")
		return err
	}
	if request == nil {
		log.Warn().Msg("DestroyNamespaceHandler: request empty error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if request.CloudCtxNs.CheckIfEmpty() {
		log.Warn().Interface("cloudCtxNs", request.CloudCtxNs).Msg("DestroyNamespaceHandler: CloudCtxNs is empty")
		return c.JSON(http.StatusBadRequest, nil)
	}
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Msg("TopologyDestroyDeploymentHandler: orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, request.CloudCtxNs)
	if authed != true {
		log.Warn().Interface("ou", ou).Interface("req", request).Msg("DestroyNamespaceHandler: IsOrgCloudCtxNsAuthorized error")
		return c.JSON(http.StatusForbidden, nil)
	}
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("DestroyNamespaceHandler: IsOrgCloudCtxNsAuthorized error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return DestroyNamespaceCluster(c, request.CloudCtxNs)
}
