package hestia_iris_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func ReadOrgRoutesRequestHandler(c echo.Context) error {
	request := new(ReadOrgRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Read(c)
}

type ReadOrgRoutesRequest struct {
}

func (r *ReadOrgRoutesRequest) Read(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	_, err := iris_models.SelectOrgRoutes(context.Background(), ou.OrgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

func ReadOrgGroupRoutesRequestHandler(c echo.Context) error {
	request := new(ReadOrgRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Read(c)
}

type ReadOrgGroupRoutesRequest struct {
}

func (r *ReadOrgGroupRoutesRequest) Read(c echo.Context) error {
	//ou := c.Get("orgUser").(org_users.OrgUser)
	return c.JSON(http.StatusOK, nil)
}
