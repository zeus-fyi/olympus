package hestia_quicknode_dashboard

import (
	"context"
	"errors"

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

type PlanUsageDetails struct {
	PlanName     string                        `json:"planName"`
	ComputeUsage *iris_usage_meters.UsageMeter `json:"computeUsage,omitempty"`
	TableUsage   iris_models.TableUsage        `json:"tableUsage"`
}

func GetUserPlanInfo(ctx context.Context, ou org_users.OrgUser, sp string) (PlanUsageDetails, error) {
	usage, err := iris_redis.IrisRedisClient.GetPlanUsageInfo(context.Background(), ou.OrgID)
	if err != nil {
		log.Err(err).Interface("usage", usage).Msg("GetPlanUsageInfo error")
		return PlanUsageDetails{}, err
	}
	switch sp {
	case "enterprise":
	case "performance":
	case "standard":
	case "lite":
	case "test":
	default:
		return PlanUsageDetails{}, errors.New("invalid service plan")
	}
	tc, err := iris_models.OrgEndpointsAndGroupTablesCount(context.Background(), ou.OrgID)
	if err != nil {
		log.Err(err).Msg("GetUserPlanInfo: OrgEndpointsAndGroupTablesCount")
		return PlanUsageDetails{}, err
	}
	usageInfo := PlanUsageDetails{
		PlanName:     sp,
		ComputeUsage: usage,
	}
	if tc != nil {
		usageInfo.TableUsage = *tc
	}
	return usageInfo, nil
}
