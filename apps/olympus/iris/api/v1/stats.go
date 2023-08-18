package v1_iris

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
)

type RateRequest struct {
}

func RateRequestHandler(c echo.Context) error {
	request := new(RateRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetUsageRate(c)
}

func (r *RateRequest) GetUsageRate(c echo.Context) error {
	ouc := c.Get("orgUser")
	ou, ok := ouc.(org_users.OrgUser)
	if !ok {
		ou = org_users.OrgUser{}
	}
	plan := ""
	sp, ok := c.Get("servicePlan").(string)
	if ok {
		plan = sp
	}
	if plan == "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "no service plan found"})
	}
	ur, err := iris_redis.IrisRedisClient.GetUsageRates(context.Background(), ou.OrgID, nil)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("RateRequest: RateRequestHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, ur)
}
