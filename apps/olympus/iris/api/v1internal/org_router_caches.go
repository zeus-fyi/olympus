package v1internal_iris

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	iris_serverless "github.com/zeus-fyi/olympus/pkg/iris/serverless"
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
	if len(orgID) == 0 {
		log.Warn().Msg("orgID is empty")
		return c.JSON(http.StatusBadRequest, nil)
	}
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

func InternalRestoreCacheForAllOrgsHandler(c echo.Context) error {
	request := new(RefreshOrgRoutingTableRequest)
	if err := c.Bind(&request); err != nil {
		log.Err(err)
		return err
	}
	return request.RefreshAllOrgRoutingTables(c)
}

func InternalRefreshServerlessTablesHandler(c echo.Context) error {
	request := new(RefreshOrgRoutingTableRequest)
	if err := c.Bind(&request); err != nil {
		log.Err(err)
		return err
	}
	return request.RefreshServerlessTables(c)
}

func (p *RefreshOrgRoutingTableRequest) RefreshServerlessTables(c echo.Context) error {
	err := iris_serverless.IrisPlatformServicesWorker.ExecuteIrisServerlessResyncWorkflow(context.Background())
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

func InternalRefreshOrgGroupRoutingTableHandler(c echo.Context) error {
	request := new(RefreshOrgRoutingTableRequest)
	if err := c.Bind(&request); err != nil {
		log.Err(err)
		return err
	}
	return request.RefreshOrgGroupRoutingTable(c)
}

func (p *RefreshOrgRoutingTableRequest) RefreshOrgGroupRoutingTable(c echo.Context) error {
	orgID := c.Param("orgID")
	if len(orgID) == 0 {
		log.Warn().Msg("DeleteOrgGroupRoutingTable: orgID is empty")
		return c.JSON(http.StatusBadRequest, nil)
	}
	oi, err := strconv.Atoi(orgID)
	if err != nil {
		log.Err(err).Msg("DeleteOrgGroupRoutingTable: strconv.Atoi(orgID)")
		return c.JSON(http.StatusBadRequest, err)
	}
	groupName := c.Param("groupName")
	if len(groupName) == 0 {
		log.Warn().Msg("DeleteOrgGroupRoutingTable: groupName is empty")
		return c.JSON(http.StatusBadRequest, nil)
	}
	err = iris_api_requests.IrisCacheWorker.ExecuteIrisCacheUpdateOrAddOrgGroupRoutingTableWorkflow(context.Background(), oi, groupName)
	if err != nil {
		log.Err(err).Msg("IrisCacheWorker.ExecuteIrisCacheUpdateOrAddOrgGroupRoutingTableWorkflow")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
