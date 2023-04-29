package replace_topology

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func TopologyDeploymentReplaceHandler(c echo.Context) error {
	request := new(TopologyReplaceRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	// from auth lookup
	// validate context kns
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	topID := c.FormValue("topologyID")
	topologyID := string_utils.IntStringParser(topID)
	request.TopologyID = topologyID
	cloudCtxNs := zeus.ConvertCloudCtxNsFormToType(c)
	request.CloudCtxNs = cloudCtxNs
	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, cloudCtxNs)
	if authed != true {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return request.ReplaceTopology(c)
}

func DeployUIClusterUpdateRequestHandler(c echo.Context) error {
	request := new(DeployClusterUpdateRequestUI)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.TopologyUpdateRequestUI(c)
}
