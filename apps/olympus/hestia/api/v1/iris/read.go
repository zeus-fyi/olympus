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

type OrgRoutesResponse struct {
	Routes []string `json:"routes"`
}

func (r *ReadOrgRoutesRequest) Read(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	routes, err := iris_models.SelectOrgRoutes(context.Background(), ou.OrgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp := OrgRoutesResponse{
		Routes: make([]string, len(routes)),
	}
	for i, route := range routes {
		resp.Routes[i] = route.RoutePath
	}
	return c.JSON(http.StatusOK, resp)
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

type OrgGroupRoutesResponse struct {
	GroupName string   `json:"groupName"`
	Routes    []string `json:"routes"`
}

func (r *ReadOrgGroupRoutesRequest) ReadGroup(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	groupName := c.Param("groupName")
	groupedRoutes, err := iris_models.SelectOrgRoutesByOrgAndGroupName(context.Background(), ou.OrgID, groupName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	routes := groupedRoutes.Map[ou.OrgID][groupName]
	resp := OrgGroupRoutesResponse{
		GroupName: groupName,
		Routes:    routes,
	}
	return c.JSON(http.StatusOK, resp)
}

func ReadOrgGroupsRoutesRequestHandler(c echo.Context) error {
	request := new(ReadOrgRoutingGroupsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ReadGroups(c)
}

type ReadOrgRoutingGroupsRequest struct {
}

type OrgGroupsRoutesResponse struct {
	Map    map[string][]string `json:"orgGroupsRoutes"`
	Routes []string            `json:"routes,omitempty"`
}

func (r *ReadOrgRoutingGroupsRequest) ReadGroups(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	groupedRoutes, err := iris_models.SelectAllOrgRoutesByOrg(context.Background(), ou.OrgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	routeGroups := groupedRoutes.Map[ou.OrgID]
	resp := OrgGroupsRoutesResponse{
		Map: routeGroups,
	}
	return c.JSON(http.StatusOK, resp)
}

func ReadAllOrgGroupsAndEndpointsRequestHandler(c echo.Context) error {
	request := new(ReadOrgRoutingGroupsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ReadAllOrgGroupsAndEndpoints(c)
}

func (r *ReadOrgRoutingGroupsRequest) ReadAllOrgGroupsAndEndpoints(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	groupedRoutes, err := iris_models.SelectAllEndpointsAndOrgGroupRoutesByOrg(context.Background(), ou.OrgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp := OrgGroupsRoutesResponse{
		Map:    groupedRoutes.Map,
		Routes: groupedRoutes.Routes,
	}
	return c.JSON(http.StatusOK, resp)
}
