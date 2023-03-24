package poseidon_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	client_consts "github.com/zeus-fyi/zeus/cookbooks/ethereum/beacons/constants"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

func (d *PoseidonSyncActivities) RestartBeaconPod(ctx context.Context, clientName string, cctx zeus_common_types.CloudCtxNs) error {
	restartClient := d.BeaconActionsClient
	restartClient.BeaconKnsReq.CloudCtxNs = cctx
	if client_consts.IsConsensusClient(clientName) {
		resp, rerr := restartClient.RestartConsensusClientPods(ctx)
		if rerr != nil {
			log.Ctx(ctx).Err(rerr).Msg("PoseidonSyncActivities: RestartBeaconPod Consensus Client")
			return rerr
		}
		log.Info().Msg(string(resp))
	}
	if client_consts.IsExecClient(clientName) {
		resp, rerr := restartClient.RestartExecClientPods(ctx)
		if rerr != nil {
			log.Ctx(ctx).Err(rerr).Msg("PoseidonSyncActivities: RestartBeaconPod Exec Client")
			return rerr
		}
		log.Info().Msg(string(resp))
	}
	return nil
}
