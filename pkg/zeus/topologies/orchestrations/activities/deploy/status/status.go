package deployment_status

import (
	"context"

	"github.com/rs/zerolog/log"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/auth"
)

const updateDeployStatusRoute = "/v1/internal/deploy/status"

type TopologyActivityDeploymentStatusActivity struct {
	Host string
	topology_deployment_status.Status
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *TopologyActivityDeploymentStatusActivity) GetActivities() ActivitiesSlice {
	return []interface{}{d.PostStatusUpdate}
}

func (d *TopologyActivityDeploymentStatusActivity) PostStatusUpdate(ctx context.Context, params topology_deployment_status.Status) error {
	_, err := api_auth_temporal.ZeusClient.R().
		SetBody(params).
		Post(zeus_client.InternalDeployStatusUpdatePath)
	if err != nil {
		log.Err(err).Interface("path", zeus_client.InternalDeployStatusUpdatePath).Msg("TopologyActivityDeploymentStatusActivity")
		return err
	}
	return err
}
