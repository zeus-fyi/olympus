package config_maps

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
)

func ConfigMapActionRequestRoutes(e *echo.Group) *echo.Group {
	e.POST("/configmaps", ConfigMapActionRequestHandler)
	return e
}

func ConfigMapActionRequestHandler(c echo.Context) error {
	request := new(ConfigMapActionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, request.CloudCtxNs)
	if authed != true {
		return c.JSON(http.StatusInternalServerError, err)
	}
	switch request.Action {
	case "key-swap":
		return request.KeySwap(c)
	case "set-or-create-from-key":
		return request.KeySetOrCreateFromExistingKey(c)
	default:
		return c.JSON(http.StatusBadRequest, nil)
	}
}
