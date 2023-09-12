package v1internal_iris

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
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

type DeleteOrgRoutingTableRequest struct {
}

func InternalDeleteOrgRoutingTableRequestHandler(c echo.Context) error {
	request := new(DeleteOrgRoutingTableRequest)
	if err := c.Bind(&request); err != nil {
		log.Err(err)
		return err
	}
	return request.DeleteOrgRoutingTables(c)
}

func (p *DeleteOrgRoutingTableRequest) DeleteOrgRoutingTables(c echo.Context) error {
	orgID := c.Param("orgID")
	if len(orgID) == 0 {
		log.Warn().Msg("DeleteOrgRoutingTables: orgID is empty")
		return c.JSON(http.StatusBadRequest, nil)
	}
	oi, err := strconv.Atoi(orgID)
	if err != nil {
		log.Err(err).Interface("orgID", orgID).Msg("DeleteOrgRoutingTables: strconv.Atoi(orgID)")
		return c.JSON(http.StatusBadRequest, err)
	}
	err = iris_api_requests.IrisCacheWorker.ExecuteIrisCacheDeleteOrgRoutingGroupTablesWorkflow(context.Background(), oi)
	if err != nil {
		log.Err(err).Msg("DeleteOrgRoutingTables: IrisCacheWorker.ExecuteIrisCacheDeleteOrgRoutingGroupTablesWorkflow")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}

func InternalDeleteOrgGroupRoutingTableRequestHandler(c echo.Context) error {
	request := new(DeleteOrgRoutingTableRequest)
	if err := c.Bind(&request); err != nil {
		log.Err(err)
		return err
	}
	return request.DeleteOrgGroupRoutingTable(c)
}

func InternalDeleteQnOrgAuthCacheHandler(c echo.Context) error {
	request := new(DeleteOrgRoutingTableRequest)
	if err := c.Bind(&request); err != nil {
		log.Err(err)
		return err
	}
	return request.DeleteQnOrgAuthCache(c)
}
func InternalDeleteSessionAuthCacheHandler(c echo.Context) error {
	request := new(DeleteOrgRoutingTableRequest)
	if err := c.Bind(&request); err != nil {
		log.Err(err)
		return err
	}
	return request.DeleteSessionIDAuthCache(c)
}

func (p *DeleteOrgRoutingTableRequest) DeleteSessionIDAuthCache(c echo.Context) error {
	sessionID := c.Param("sessionID")
	if len(sessionID) == 0 {
		log.Warn().Msg("DeleteQnOrgAuthCache: orgID is empty")
		return c.JSON(http.StatusBadRequest, nil)
	}
	err := iris_redis.IrisRedisClient.DeleteAuthCache(context.Background(), sessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, "ok")
}

func (p *DeleteOrgRoutingTableRequest) DeleteQnOrgAuthCache(c echo.Context) error {
	qnID := c.Param("qnID")
	if len(qnID) == 0 {
		log.Warn().Msg("DeleteQnOrgAuthCache: orgID is empty")
		return c.JSON(http.StatusBadRequest, nil)
	}
	key := read_keys.NewKeyReader()
	keysFound, err := key.GetUserKeysToServices(context.Background(), qnID)
	if err != nil {
		log.Err(err).Msg("DeleteQnOrgAuthCache: QueryUserAuthedServices error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	for k, _ := range keysFound {
		err = iris_redis.IrisRedisClient.DeleteAuthCache(context.Background(), k)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	return c.JSON(http.StatusOK, "ok")
}

func (p *DeleteOrgRoutingTableRequest) DeleteOrgGroupRoutingTable(c echo.Context) error {
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
	err = iris_api_requests.IrisCacheWorker.ExecuteIrisCacheDeleteOrgGroupRoutingTableWorkflow(context.Background(), oi, groupName)
	if err != nil {
		log.Err(err).Msg("DeleteOrgGroupRoutingTable: IrisCacheWorker.ExecuteIrisCacheDeleteOrgGroupRoutingTableWorkflow")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
