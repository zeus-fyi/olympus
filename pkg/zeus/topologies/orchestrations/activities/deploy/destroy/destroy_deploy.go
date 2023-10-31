package destroy_deploy_activities

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type DestroyDeployTopologyActivities struct {
	base_deploy_params.TopologyWorkflowRequest
}
type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *DestroyDeployTopologyActivities) GetActivities() ActivitiesSlice {
	return []interface{}{
		d.DestroyNamespace,
		d.DestroyDeployStatefulSet,
		d.DestroyDeployDeployment,
		d.DestroyDeployConfigMap,
		d.DestroyDeployService,
		d.DestroyDeployIngress,
		d.DestroyDeployServiceMonitor,
		d.DestroyJob,
		d.DestroyCronJob,
		d.DeletePod,
	}
}

func (d *DestroyDeployTopologyActivities) postDestroyDeployTarget(target string, params base_request.InternalDeploymentActionRequest) error {
	u := d.GetDestroyDeployURL(target)
	client := resty.New()
	client.SetBaseURL(u.Host)
	resp, err := client.R().
		SetAuthToken(api_auth_temporal.Bearer).
		SetBody(params).
		Post(u.Path)
	if err != nil {
		log.Err(err).Interface("path", u.Path).Msg("DestroyDeployTopologyActivities: postDestroyDeployTarget failed")
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		log.Err(err).Interface("path", u.Path).Interface("statusCode", resp.StatusCode()).Msg("DeployTopologyActivities: postDestroyDeployTarget failed bad status code")
		return errors.New(fmt.Sprintf("DeployTopologyActivities: postDestroyDeployTarget failed statusCode: %d", resp.StatusCode()))
	}
	return err
}

func (d *DestroyDeployTopologyActivities) GetDestroyDeployURL(target string) url.URL {
	return d.GetURL(zeus_endpoints.InternalDestroyDeployPath, target)
}

func (d *DestroyDeployTopologyActivities) DeletePod(ctx context.Context, podName string, cctx zeus_common_types.CloudCtxNs) error {
	err := zeus.K8Util.DeleteFirstPodLike(ctx, cctx, podName, nil, nil)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("DestroyDeployTopologyActivities: PodsDeleteRequest: DeleteFirstPodLike")
		return err
	}
	return nil
}
