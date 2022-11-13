package deploy_topology_activities

import (
	"net/url"

	"github.com/rs/zerolog/log"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/actions/base_request"
)

type DeployTopologyActivities struct {
	base_deploy_params.TopologyWorkflowRequest
}
type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *DeployTopologyActivities) GetActivities() ActivitiesSlice {
	return []interface{}{
		d.CreateNamespace,
		d.DeployDeployment, d.DeployStatefulSet,
		d.DeployConfigMap,
		d.DeployService, d.DeployIngress,
	}
}

func (d *DeployTopologyActivities) postDeployTarget(target string, params base_request.InternalDeploymentActionRequest) error {
	u := d.GetDeployURL(target)
	_, err := api_auth_temporal.ZeusClient.R().
		SetBody(params).
		Post(u.Path)

	if err != nil {
		log.Err(err).Interface("path", u.Path).Msg("DeployTopologyActivities: postDeployTarget failed")
		return err
	}
	return err
}

func (d *DeployTopologyActivities) GetDeployURL(target string) url.URL {
	return d.GetURL(zeus_client.InternalDeployPath, target)
}
