package deploy_topology_activities_create_setup

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

func (c *CreateSetupTopologyActivities) GetClusterAuthCtx(ctx context.Context, ou org_users.OrgUser, params zeus_common_types.CloudCtxNs) (*authorized_clusters.K8sClusterConfig, error) {
	kc, err := authorized_clusters.SelectAuthedClusterByRouteOnlyAndOrgID(ctx, ou, params)
	if err != nil {
		log.Err(err).Interface("ou", ou).Interface("params", params).Msg("CreateSetupTopologyActivities: GetClusterAuthCtx: SelectAuthedClusterByRouteOnlyAndOrgID error")
		return nil, err
	}
	return kc, nil
}

func (c *CreateSetupTopologyActivities) GetClusterAuthCtxFromID(ctx context.Context, ou org_users.OrgUser, clusterCfgStrID string) (*authorized_clusters.K8sClusterConfig, error) {
	kc, err := authorized_clusters.SelectAuthedClusterByIDAndOrgID(ctx, ou, clusterCfgStrID)
	if err != nil {
		log.Err(err).Interface("ou", ou).Interface("clusterCfgStrID", clusterCfgStrID).Msg("CreateSetupTopologyActivities: GetClusterAuthCtxFromID: error")
		return nil, err
	}
	return kc, nil
}
