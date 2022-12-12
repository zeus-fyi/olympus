package beacon_actions

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	client_consts "github.com/zeus-fyi/olympus/cookbooks/ethereum/beacons/constants"
	zeus_configmap_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/config_maps"
)

func (b *BeaconActionsClient) PauseClient(ctx context.Context, cmName, clientName string) ([]byte, error) {
	cmr := zeus_configmap_reqs.ConfigMapActionRequest{
		TopologyDeployRequest: b.BeaconKnsReq,
		Action:                zeus_configmap_reqs.SetOrCreateKeyFromExisting,
		ConfigMapName:         cmName,
		Keys: zeus_configmap_reqs.KeySwap{
			KeyOne: "pause.sh",
			KeyTwo: "start.sh",
		},
		FilterOpts: nil,
	}
	respCm, err := b.SetOrCreateKeyFromConfigMapKey(ctx, cmr)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("configMap", respCm).Msg("PauseConsensusClient: SetOrCreateKeyFromConfigMapKey")
		return nil, err
	}
	if client_consts.IsConsensusClient(clientName) {
		b.ConsensusClient = clientName
		resp, cerr := b.RestartConsensusClientPods(ctx, basePar)
		if cerr != nil {
			log.Ctx(ctx).Err(cerr).Msg("PauseConsensusClient: RestartConsensusClientPods")
			return nil, cerr
		}
		return resp, cerr
	}
	if client_consts.IsExecClient(clientName) {
		b.ExecClient = clientName
		resp, cerr := b.RestartExecClientPods(ctx, basePar)
		if cerr != nil {
			log.Ctx(ctx).Err(cerr).Msg("PauseConsensusClient: RestartExecClientPods")
			return nil, cerr
		}
		return resp, cerr
	}
	return nil, errors.New("invalid consensus exec client supplied")
}
