package hestia_iris_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris/models/bases/autogen"
)

func CreateOrgRoutesRequestHandler(c echo.Context) error {
	request := new(CreateOrgRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Create(c)
}

type CreateOrgRoutesRequest struct {
	Routes []string `json:"routes"`
}

func (r *CreateOrgRoutesRequest) Create(c echo.Context) error {
	or := make([]iris_autogen_bases.OrgRoutes, len(r.Routes))
	for i, route := range r.Routes {
		or[i] = iris_autogen_bases.OrgRoutes{
			RoutePath: route,
		}
	}
	err := iris_models.InsertOrgRoutes(context.Background(), or)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

func CreateOrgGroupRoutesRequestHandler(c echo.Context) error {
	request := new(CreateOrgRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Create(c)
}

type CreateOrUpdateOrgGroupRoutesRequest struct {
	GroupName string   `json:"groupName"`
	Routes    []string `json:"routes"`
}

func (r *CreateOrUpdateOrgGroupRoutesRequest) Create(c echo.Context) error {
	or := make([]iris_autogen_bases.OrgRoutes, len(r.Routes))
	for i, route := range r.Routes {
		or[i] = iris_autogen_bases.OrgRoutes{
			RoutePath: route,
		}
	}
	err := iris_models.InsertOrgRoutes(context.Background(), or)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}
