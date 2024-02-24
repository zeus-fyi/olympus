package deploy_updates

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func DeployUIClusterRolloutRestartRequestHandler(c echo.Context) error {
	return TopologyRolloutRestartRequestUI(c)
}

func TopologyRolloutRestartRequestUI(c echo.Context) error {
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Msg("TopologyRolloutRestartRequestUI, orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	cctxID := c.Request().Header.Get("CloudCtxNsID")
	cID, err := strconv.Atoi(cctxID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	authed, cctx, err := read_topology.IsOrgCloudCtxNsAuthorizedFromID(ctx, ou.OrgID, cID)
	if authed != true {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	if cctx.CheckIfEmpty() {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return zeus.ExecuteDeployRolloutRestartWorkflow(c, ctx, ou, cID, cctx)
}
func DeployApiClusterRolloutRestartRequestHandler(c echo.Context) error {
	return TopologyRolloutRestartRequestApi(c)
}

func TopologyRolloutRestartRequestApi(c echo.Context) error {
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Msg("TopologyRolloutRestartRequestUI, orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	cctxID := c.Request().Header.Get("CloudCtxNsID")
	cID, err := strconv.Atoi(cctxID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	authed, cctx, err := read_topology.IsOrgCloudCtxNsAuthorizedFromID(ctx, ou.OrgID, cID)
	if authed != true {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	if cctx.CheckIfEmpty() {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return zeus.ExecuteDeployRolloutRestartWorkflow(c, ctx, ou, cID, cctx)
}
