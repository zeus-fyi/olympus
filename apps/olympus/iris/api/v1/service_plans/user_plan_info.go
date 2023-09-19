package iris_service_plans

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
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

const (
	PlanUsageDetailsRoute    = "/plan/usage"
	TableMetricsDetailsRoute = "/table/:groupName/metrics"
)

type PlanUsageDetailsRequest struct {
}

type PlanUsageDetailsResponse struct {
	PlanName     string                                `json:"planName"`
	ComputeUsage *iris_usage_meters.UsageMeter         `json:"computeUsage,omitempty"`
	TableUsage   iris_models.TableUsageAndUserSettings `json:"tableUsage"`
}

func PlanUsageDetailsRequestHandler(c echo.Context) error {
	log.Info().Msg("Iris: PlanUsageDetailsRequest")
	request := new(PlanUsageDetailsRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("Iris: PlanUsageDetailsRequest")
		return err
	}
	return request.GetUserPlanInfo(c)
}

func (p *PlanUsageDetailsRequest) GetUserPlanInfo(c echo.Context) error {
	planName, ok := c.Get("servicePlan").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ou := org_users.OrgUser{}
	ouc := c.Get("orgUser")
	if ouc != nil {
		ouser, aok := ouc.(org_users.OrgUser)
		if aok {
			ou = ouser
		}
	}
	usage, err := iris_redis.IrisRedisClient.GetPlanUsageInfo(context.Background(), ou.OrgID)
	if err != nil {
		log.Err(err).Interface("usage", usage).Msg("GetPlanUsageInfo error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	switch planName {
	case "enterprise":
		planName = "Enterprise"
	case "performance":
		planName = "Performance"
	case "standard":
		planName = "Standard"
	case "lite":
		planName = "Lite"
	case "test":
		planName = "Test"
	default:
		return c.JSON(http.StatusInternalServerError, nil)
	}

	usage.MonthlyBudgetZU = float64(iris_redis.GetMonthlyPlanBudgetZU(planName))
	usage.RateLimit = float64(iris_redis.GetMonthlyPlanBudgetThroughputZU(planName))
	usage.GetMonthlyUsageZUM()
	usage.GetMonthlyBudgetZUM()
	usage.GetCurrentRateZUk()
	usage.GetRateLimitZUk()
	tc, err := iris_models.OrgEndpointsAndGroupTablesCount(context.Background(), ou.OrgID, ou.UserID)
	if err != nil {
		log.Err(err).Msg("GetUserPlanInfo: OrgEndpointsAndGroupTablesCount")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	usageInfo := PlanUsageDetailsResponse{
		PlanName:     planName,
		ComputeUsage: usage,
	}
	if tc != nil {
		err = tc.SetMaxTableCountByPlan(planName)
		if err != nil {
			log.Err(err).Msg("GetUserPlanInfo: SetMaxTableCountByPlan")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		usageInfo.TableUsage = *tc
	}
	return c.JSON(http.StatusOK, usageInfo)
}
