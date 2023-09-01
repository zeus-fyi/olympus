package iris_service_plans

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
)

func TableMetricsDetailsRequestHandler(c echo.Context) error {
	request := new(PlanUsageDetailsRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("Iris: GetTableMetrics")
		return err
	}
	return request.GetTableMetrics(c)
}

func (p *PlanUsageDetailsRequest) GetTableMetrics(c echo.Context) error {
	ou := org_users.OrgUser{}
	ouc := c.Get("orgUser")
	if ouc != nil {
		ouser, aok := ouc.(org_users.OrgUser)
		if aok {
			ou = ouser
		}
	}
	tblName := c.Param("groupName")
	usage, err := iris_redis.IrisRedisClient.GetPriorityScoresAndTdigestMetrics(context.Background(), ou.OrgID, tblName)
	if err != nil {
		log.Err(err).Interface("usage", usage).Msg("GetTableMetrics: GetPriorityScoresAndTdigestMetrics error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, usage)
}
