package poseidon_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	pg_poseidon "github.com/zeus-fyi/olympus/datastores/postgres/apps/poseidon"
	client_consts "github.com/zeus-fyi/zeus/cookbooks/ethereum/beacons/constants"
)

func (d *PoseidonSyncActivities) ScheduleDiskWipe(ctx context.Context, params pg_poseidon.DiskWipeOrchestration) error {
	err := params.ScheduleDiskWipe(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("PoseidonSyncActivities: ScheduleDiskWipe")
		return err
	}
	return err
}

func (d *PoseidonSyncActivities) RestartBeaconPod(ctx context.Context, params pg_poseidon.DiskWipeOrchestration) error {
	err := params.ScheduleDiskWipe(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("PoseidonSyncActivities: ScheduleDiskWipe")
		return err
	}
	restartClient := d.BeaconActionsClient
	restartClient.BeaconKnsReq.CloudCtxNs = params.CloudCtxNs
	if client_consts.IsConsensusClient(params.ClientName) {
		resp, rerr := restartClient.RestartConsensusClientPods(ctx)
		if rerr != nil {
			log.Ctx(ctx).Err(rerr).Msg("PoseidonSyncActivities: ScheduleDiskWipe Consensus Client")
			return rerr
		}
		log.Info().Msg(string(resp))
	}
	if client_consts.IsExecClient(params.ClientName) {
		resp, rerr := restartClient.RestartExecClientPods(ctx)
		if rerr != nil {
			log.Ctx(ctx).Err(rerr).Msg("PoseidonSyncActivities: ScheduleDiskWipe Exec Client")
			return rerr
		}
		log.Info().Msg(string(resp))
	}
	return err
}

func (d *PoseidonSyncActivities) ScheduleDiskUpload(ctx context.Context, params pg_poseidon.DiskWipeOrchestration) error {
	err := params.ScheduleDiskWipe(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("PoseidonSyncActivities: ScheduleDiskWipe")
		return err
	}
	return err
}
