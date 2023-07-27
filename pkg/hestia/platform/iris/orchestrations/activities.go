package platform_service_orchestrations

import (
	"context"
)

type HestiaPlatformActivities struct {
}

func NewHestiaPlatformActivities() HestiaPlatformActivities {
	return HestiaPlatformActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (h *HestiaPlatformActivities) GetActivities() ActivitiesSlice {
	return []interface{}{
		h.IrisPlatformSetupCacheUpdateRequest, h.UpdateDatabaseOrgRoutingTables,
	}
}

func (h *HestiaPlatformActivities) UpdateDatabaseOrgRoutingTables(ctx context.Context, pr IrisPlatformServiceRequest) error {
	//err := iris_models.InsertOrgRoutes(context.Background(), pr.Routes)
	//if err != nil {
	//}
	if pr.OrgGroupName != "" {
		// then do group routes
	}
	return nil
}

func (h *HestiaPlatformActivities) IrisPlatformSetupCacheUpdateRequest(ctx context.Context, pr IrisPlatformServiceRequest) error {
	// call iris via api & cache refresh
	return nil
}
