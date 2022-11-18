package create_or_update_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
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
	ou := c.Get("orgUser").(org_users.OrgUser)
	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, request.CloudCtxNs)
	if authed != true {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return request.DeployTopology(c)
}
