package v1internal_iris

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
)

func InternalDeleteQnOrgAuthCacheHandler(c echo.Context) error {
	request := new(DeleteOrgRoutingTableRequest)
	if err := c.Bind(&request); err != nil {
		log.Err(err)
		return err
	}
	return request.DeleteQnOrgAuthCache(c)
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
