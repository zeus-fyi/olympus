package v1internal_iris

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
)

type RefreshOrgRoutingTableRequest struct {
}

func InternalRefreshOrgRoutingTableHandler(c echo.Context) error {
	request := new(RefreshOrgRoutingTableRequest)
	if err := c.Bind(&request); err != nil {
		log.Err(err)
		return err
	}
	return request.RefreshOrgRoutingTable(c)
}
func (p *RefreshOrgRoutingTableRequest) RefreshOrgRoutingTable(c echo.Context) error {
	orgID := c.Param("orgID")
	oi, err := strconv.Atoi(orgID)
	if err != nil {
		log.Err(err).Msg("strconv.Atoi(orgID)")
		return c.JSON(http.StatusBadRequest, err)
	}
	err = iris_api_requests.IrisCacheWorker.ExecuteIrisCacheUpdateOrAddOrgRoutingTablesWorkflow(context.Background(), oi)
	if err != nil {
		log.Err(err).Msg("IrisCacheWorker.ExecuteIrisCacheUpdateOrAddOrgRoutingTablesWorkflow")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}

func (p *RefreshOrgRoutingTableRequest) RefreshAllOrgRoutingTables(c echo.Context) error {
	err := iris_api_requests.IrisCacheWorker.ExecuteIrisCacheRefreshAllOrgRoutingTablesWorkflow(context.Background())
	if err != nil {
		log.Err(err).Msg("IrisCacheWorker.ExecuteIrisCacheUpdateOrAddOrgRoutingTablesWorkflow")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
