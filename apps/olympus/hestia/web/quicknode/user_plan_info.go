package hestia_quicknode_dashboard

import (
	"context"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
)

/*
todo: get user plan info


1. user plan name
2. user plan consumption
	a. per month ZU consumption
	b. endpoints x/1000
	c. custom tables x/(plan specific)
3. per table metrics histograms
	a. eg. t-digests per metric
*/

type PlanUsage struct {
	PlanName     string                 `json:"planName"`
	MonthlyUsage float64                `json:"monthlyUsage"`
	TableUsage   iris_models.TableUsage `json:"tableUsage"`
}

// TODO, should get this on login & jwt refresh

func GetUserPlanInfo(c echo.Context) (PlanUsage, error) {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return PlanUsage{}, errors.New("failed to cast orgUser")
	}

	usage, err := iris_redis.IrisRedisClient.GetPlanUsageInfo(context.Background(), ou.OrgID)
	if err != nil {
		log.Err(err).Interface("usage", usage).Msg("GetPlanUsageInfo error")
		return PlanUsage{}, err
	}
	sp, ok := c.Get("servicePlans").(map[string]string)
	if !ok {
		log.Warn().Interface("servicePlans", sp).Msg("GetUserPlanInfo: marketplace plan not found")
		return PlanUsage{}, err
	}

	plan := "todo"
	switch plan {
	case "enterprise":
	case "performance":
	case "standard":
	case "lite":
	}

	tc, err := iris_models.OrgEndpointsAndGroupTablesCount(context.Background(), ou.OrgID)
	if err != nil {
		log.Err(err).Msg("GetUserPlanInfo: OrgEndpointsAndGroupTablesCount")
		return PlanUsage{}, err
	}
	usageInfo := PlanUsage{
		PlanName:     plan,
		MonthlyUsage: 0,
	}
	if tc != nil {
		usageInfo.TableUsage = *tc
	}
	return usageInfo, nil
}
