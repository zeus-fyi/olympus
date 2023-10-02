package hestia_web_resources

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
)

type ResourceListRequest struct {
}

func ResourceListRequestHandler(c echo.Context) error {
	request := new(ResourceListRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.List(c)
}

func (r *ResourceListRequest) List(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	ctx := context.Background()
	nl, err := hestia_compute_resources.SelectOrgResourcesNodes(ctx, ou.OrgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nl)
	}
	return c.JSON(http.StatusOK, nl)
}
