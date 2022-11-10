package destroy_deploy_activities

import (
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
)

const destroyDeployRoute = "/v1/internal/deploy/destroy"

type DestroyDeployTopologyActivities struct {
	base_deploy_params.TopologyWorkflowRequest
}
type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *DestroyDeployTopologyActivities) GetActivities() ActivitiesSlice {
	return []interface{}{
		d.DestroyNamespace,
		d.DestroyDeployStatefulSet, d.DestroyDeployDeployment,
		d.DestroyDeployConfigMap,
		d.DestroyDeployService, d.DestroyDeployIngress,
	}
}

func (d *DestroyDeployTopologyActivities) postDestroyDeployTarget(target string) error {
	u := d.GetDestroyDeployURL(target)
	client := resty.New()
	_, err := client.R().
		SetAuthToken(api_auth_temporal.Bearer).
		SetBody(d.TopologyWorkflowRequest).
		Post(u.Path)
	if err != nil {
		log.Err(err).Interface("path", u.Path).Msg("DestroyDeployTopologyActivities: postDestroyDeployTarget failed")
		return err
	}
	return err
}

func (d *DestroyDeployTopologyActivities) GetDestroyDeployURL(target string) url.URL {
	return d.GetURL(destroyDeployRoute, target)
}
