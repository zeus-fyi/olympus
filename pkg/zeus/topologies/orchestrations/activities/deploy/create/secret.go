package deploy_topology_activities

import (
	"context"

	aegis_secrets "github.com/zeus-fyi/olympus/datastores/postgres/apps/aegis"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
)

func (d *DeployTopologyActivities) CreateChoreographySecret(ctx context.Context, params base_request.InternalDeploymentActionRequest) error {
	return d.postDeployTarget("choreography/secrets", params)
}

func (d *DeployTopologyActivities) CreateSecret(ctx context.Context, params base_request.InternalDeploymentActionRequest, secretRef string) error {
	exists, err := aegis_secrets.DoesOrgSecretExistForTopology(ctx, params.OrgUser.OrgID, secretRef)
	if err != nil {
		return err
	}
	if exists {
		return d.postDeployTarget("dynamic/secrets", params)
	}
	return nil
}
