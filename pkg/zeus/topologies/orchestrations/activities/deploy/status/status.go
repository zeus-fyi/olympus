package deployment_status

import (
	"context"
	"net/url"

	"github.com/rs/zerolog/log"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	create_kns "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/kns"
	create_topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/state"
	delete_kns "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/delete/topologies/topology/kns"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

type TopologyActivityDeploymentStatusActivity struct {
	Host string
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *TopologyActivityDeploymentStatusActivity) GetActivities() ActivitiesSlice {
	return []interface{}{d.PostStatusUpdate, d.CreateOrUpdateKubeCtxNsStatus, d.DeleteKubeCtxNsStatus}
}

func (d *TopologyActivityDeploymentStatusActivity) PostStatusUpdate(ctx context.Context, status topology_deployment_status.DeployStatus) error {
	err := create_topology_deployment_status.InsertOrUpdateStatus(ctx, &status)
	if err != nil {
		log.Err(err).Interface("status", status).Msg("UpdateWorkloadStateHandler")
		return err
	}
	return err
}

func (d *TopologyActivityDeploymentStatusActivity) CreateOrUpdateKubeCtxNsStatus(ctx context.Context, topDepReq zeus_req_types.TopologyDeployRequest) error {
	err := create_kns.InsertKns(ctx, &topDepReq)
	if err != nil {
		log.Err(err).Interface("kns", topDepReq).Msg("InsertOrUpdateWorkloadKnsStateHandler")
		return err
	}
	return nil
}

func (d *TopologyActivityDeploymentStatusActivity) DeleteKubeCtxNsStatus(ctx context.Context, topDepReq zeus_req_types.TopologyDeployRequest) error {
	if topDepReq.TopologyID == 0 {
		err := delete_kns.DeleteKnsByOrgAccessAndCloudCtx(ctx, &topDepReq)
		if err != nil {
			log.Err(err).Interface("kns", topDepReq).Msg("DeleteWorkloadKnsStateHandler")
			return err
		}
	} else {
		err := delete_kns.DeleteKns(ctx, &topDepReq)
		if err != nil {
			log.Err(err).Interface("kns", topDepReq).Msg("DeleteWorkloadKnsStateHandler")
			return err
		}
	}
	return nil
}

func (d *TopologyActivityDeploymentStatusActivity) GetDeploymentStatusUpdateURL() url.URL {
	return d.GetURL()
}

var BaseURL = "https://api.zeus.fyi"

func (d *TopologyActivityDeploymentStatusActivity) GetURL() url.URL {
	if len(d.Host) <= 0 {
		d.Host = BaseURL
	}
	u := url.URL{
		Host: d.Host,
	}
	return u
}
