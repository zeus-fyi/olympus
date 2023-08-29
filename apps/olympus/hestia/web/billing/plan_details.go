package hestia_billing

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	iris_service_plans "github.com/zeus-fyi/olympus/iris/api/v1/service_plans"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

var (
	IrisApiUrl           = "https://iris.zeus.fyi"
	ArtificialTableCount = 2
)

func GetPlan(ctx context.Context, token string) (iris_service_plans.PlanUsageDetailsResponse, error) {
	planUsageDetails := iris_service_plans.PlanUsageDetailsResponse{}
	rc := resty_base.GetBaseRestyClient(IrisApiUrl, token)
	refreshEndpoint := fmt.Sprintf("/v1/plan/usage")
	resp, err := rc.R().SetResult(&planUsageDetails).Get(refreshEndpoint)
	if err != nil {
		log.Err(err).Msg("GetPlan: IrisPlatformSetupCacheUpdateRequest")
		return planUsageDetails, err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Msg("GetPlan: IrisPlatformSetupCacheUpdateRequest")
		return planUsageDetails, err
	}
	planUsageDetails.TableUsage.TableCount += ArtificialTableCount
	return planUsageDetails, err
}
