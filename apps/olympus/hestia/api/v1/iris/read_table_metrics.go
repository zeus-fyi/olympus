package hestia_iris_v1_routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type ReadMetricsRequest struct{}

func ReadTableMetrics(c echo.Context) error {
	request := new(ReadMetricsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ReadTableStats(c)
}

type TableMetricsResponse struct {
}

// TODO: implement

func GetTableMetrics(ou org_users.OrgUser) (TableMetricsResponse, error) {
	return TableMetricsResponse{}, nil
}

func (r *ReadMetricsRequest) ReadTableStats(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusBadRequest, "no account found")
	}
	resp, err := GetTableMetrics(ou)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, resp)
}
