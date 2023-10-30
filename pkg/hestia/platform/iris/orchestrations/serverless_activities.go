package platform_service_orchestrations

import (
	"context"
)

/*
orchestrations
	needs to auto-populate the serverless routing table
	need to add garbage collection orchestration
	auto-scaling up/down
		needs to trigger based on threshold low anvil servers in router


TODO: need to add alert
*/

func (h *HestiaPlatformActivities) ResyncServerlessRoutes(ctx context.Context) error {

	// TODO
	return nil
}
