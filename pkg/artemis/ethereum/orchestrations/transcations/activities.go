package artemis_ethereum_transcations

import (
	"net/url"

	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
)

type ArtemisEthereumBroadcastTxActivities struct {
	base_deploy_params.TopologyWorkflowRequest
}
type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *ArtemisEthereumBroadcastTxActivities) GetActivities() ActivitiesSlice {
	return []interface{}{}
}

func (d *ArtemisEthereumBroadcastTxActivities) postDeployTarget(target string, params interface{}) error {
	return nil
}

func (d *ArtemisEthereumBroadcastTxActivities) GetDeployURL(target string) url.URL {
	return d.GetURL(zeus_endpoints.InternalDeployPath, target)
}
