package iris_serverless

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
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

const (
	AnvilServerlessRoutingTable = "anvil"
)

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
		i.ResyncServerlessRoutes, i.FetchLatestServerlessRoutes,
	}
	return actSlice
}

// AddRoutesToServerlessRoutingTable if no routes are provided, it will just refresh any existing routes back into availability

func (i *IrisPlatformActivities) ResyncServerlessRoutes(ctx context.Context, routes []iris_models.RouteInfo) error {
	err := iris_redis.IrisRedisClient.AddRoutesToServerlessRoutingTable(ctx, AnvilServerlessRoutingTable, routes)
	if err != nil {
		log.Err(err).Msg("ResyncServerlessRoutes: AddRoutesToServerlessRoutingTable")
		return err
	}
	return nil
}

/*
Services
A records
"Normal" (not headless) Services are assigned a DNS A record for a name of the form my-svc.my-namespace.svc.cluster.local. This resolves to the cluster IP of the Service.

"Headless" (without a cluster IP) Services are also assigned a DNS A record for a name of the form my-svc.my-namespace.svc.cluster.local.
Unlike normal Services, this resolves to the set of IPs of the pods selected by the Service. Clients are expected to consume the set or else use standard round-robin selection from the set.

SRV records
SRV Records are created for named ports that are part of normal or Headless Services. For each named port,
the SRV record would have the form _my-port-name._my-port-protocol.my-svc.my-namespace.svc.cluster.local.
For a regular service, this resolves to the port number and the CNAME: my-svc.my-namespace.svc.cluster.local.
For a headless service, this resolves to multiple answers, one for each pod that is backing the service, and contains the port number
and a CNAME of the pod of the form auto-generated-name.my-svc.my-namespace.svc.cluster.local.

*/
// <service-name>.<namespace>.svc.cluster.local
// http://anvil-0.anvil.anvil-serverless-4d383226.svc.cluster.local
//internalLB = "http://anvil.anvil-serverless-4d383226.svc.cluster.local/v2/internal/router"

func (i *IrisPlatformActivities) FetchLatestServerlessRoutes(ctx context.Context) ([]iris_models.RouteInfo, error) {
	count := 2
	routes := make([]iris_models.RouteInfo, count)
	for j := 0; j < count; j++ {
		routes[j].RoutePath = fmt.Sprintf("http://anvil-%d.anvil.anvil-serverless-4d383226.svc.cluster.local", j)
	}
	return routes, nil
}