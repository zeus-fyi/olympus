package artemis_ethereum_transcations

import (
	"context"
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
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

func (d *ArtemisEthereumBroadcastTxActivities) SendEther(ctx context.Context, payload web3_actions.SendEtherPayload) error {

	send, err := ArtemisEthereumBroadcastTxClient.Send(ctx, payload)
	if err != nil {
		log.Err(err).Interface("tx", send).Interface("payload", payload).Msg("ArtemisEthereumBroadcastTxActivities: Send failed")
		return err
	}
	return nil
}

func (d *ArtemisEthereumBroadcastTxActivities) GetDeployURL(target string) url.URL {
	return d.GetURL(zeus_endpoints.InternalDeployPath, target)
}
