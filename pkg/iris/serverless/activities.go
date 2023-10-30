package iris_serverless

import (
	"context"

	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
)

/*
orchestrations
	needs to auto-populate the serverless routing table
	need to add garbage collection orchestration
	auto-scaling up/down
		needs to trigger based on threshold low anvil servers in router


TODO: need to add alert
*/

type IrisPlatformActivities struct {
	kronos_helix.KronosActivities
}

func NewIrisPlatformActivities() IrisPlatformActivities {
	return IrisPlatformActivities{
		KronosActivities: kronos_helix.NewKronosActivities(),
	}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (i *IrisPlatformActivities) GetActivities() ActivitiesSlice {
	actSlice := []interface{}{
		i.ResyncServerlessRoutes,
	}
	return actSlice
}

func (i *IrisPlatformActivities) ResyncServerlessRoutes(ctx context.Context) error {
	//iris_redis.AddRoutesToServerlessRoutingTable
	// TODO
	return nil
}
