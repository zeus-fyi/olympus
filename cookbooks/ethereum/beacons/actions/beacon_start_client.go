package beacon_actions

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	client_consts "github.com/zeus-fyi/olympus/cookbooks/ethereum/beacons/constants"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	zeus_configmap_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/config_maps"
	zeus_pods_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/pods"
)

func (b *BeaconActionsClient) StartConsensusClient(ctx context.Context) ([]byte, error) {
	return b.StartClient(ctx, "cm-lighthouse", client_consts.Lighthouse)
}

func (b *BeaconActionsClient) StartExecClient(ctx context.Context) ([]byte, error) {
	return b.StartClient(ctx, "cm-geth", client_consts.Geth)
}

func (b *BeaconActionsClient) StartClient(ctx context.Context, cmName, clientName string) ([]byte, error) {
	fo := string_utils.FilterOpts{
		DoesNotStartWithThese: nil,
		StartsWithThese:       nil,
		StartsWith:            "",
		Contains:              clientName,
		DoesNotInclude:        nil,
	}
	par := zeus_pods_reqs.PodActionRequest{
		TopologyDeployRequest: b.BeaconKnsReq,
		Action:                zeus_pods_reqs.DeleteAllPods,
		PodName:               "",
		ContainerName:         "",
		FilterOpts:            &fo,
		ClientReq:             nil,
		LogOpts:               nil,
		DeleteOpts:            nil,
	}
	cmReq := zeus_configmap_reqs.ConfigMapActionRequest{
		TopologyDeployRequest: b.BeaconKnsReq,
		Action:                zeus_configmap_reqs.SetOrCreateKeyFromExisting,
		ConfigMapName:         cmName,
		Keys: zeus_configmap_reqs.KeySwap{
			KeyOne: fmt.Sprintf("%s.sh", clientName),
			KeyTwo: "start.sh",
		},
		FilterOpts: nil,
	}
	respCm, err := b.SetOrCreateKeyFromConfigMapKey(ctx, cmReq)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("configMap", respCm).Msg("PauseConsensusClient: SetOrCreateKeyFromConfigMapKey")
		return nil, err
	}
	if client_consts.IsConsensusClient(clientName) {
		b.ConsensusClient = clientName
		resp, cerr := b.RestartConsensusClientPods(ctx, par)
		if cerr != nil {
			log.Ctx(ctx).Err(cerr).Msg("PauseConsensusClient: RestartConsensusClientPods")
			return nil, cerr
		}
		return resp, cerr
	}
	if client_consts.IsExecClient(clientName) {
		b.ExecClient = clientName
		resp, cerr := b.RestartExecClientPods(ctx, par)
		if cerr != nil {
			log.Ctx(ctx).Err(cerr).Msg("PauseConsensusClient: RestartExecClientPods")
			return nil, cerr
		}
		return resp, cerr
	}
	return nil, errors.New("invalid consensus exec client supplied")
}
