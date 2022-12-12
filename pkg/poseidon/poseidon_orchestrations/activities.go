package poseidon_orchestrations

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	beacon_actions "github.com/zeus-fyi/olympus/cookbooks/ethereum/beacons/actions"
	athena_client "github.com/zeus-fyi/olympus/pkg/athena/client"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_buckets"
)

type PoseidonSyncActivities struct {
	beacon_actions.BeaconActionsClient
	athena_client.AthenaClient
}

func NewPoseidonSyncActivity(ba beacon_actions.BeaconActionsClient, ac athena_client.AthenaClient) PoseidonSyncActivities {
	return PoseidonSyncActivities{ba, ac}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

var PoseidonSyncActivitiesOrchestrator PoseidonSyncActivities

func (d *PoseidonSyncActivities) GetActivities() ActivitiesSlice {
	return []interface{}{d.PauseExecClient, d.PauseConsensusClient, d.Resume, d.IsExecClientSynced, d.IsConsensusClientSynced, d.RsyncExecBucket, d.RsyncConsensusBucket}
}

func (d *PoseidonSyncActivities) PauseExecClient(ctx context.Context) error {
	cmName := fmt.Sprintf("cm-%s", d.ExecClient)
	_, err := d.BeaconActionsClient.PauseClient(ctx, cmName, d.ExecClient)
	return err
}

func (d *PoseidonSyncActivities) PauseConsensusClient(ctx context.Context) error {
	cmName := fmt.Sprintf("cm-%s", d.ConsensusClient)
	_, err := d.BeaconActionsClient.PauseClient(ctx, cmName, d.ConsensusClient)
	return err
}

func (d *PoseidonSyncActivities) ResumeExecClient(ctx context.Context) error {
	cmName := fmt.Sprintf("cm-%s", d.ExecClient)
	_, err := d.BeaconActionsClient.StartClient(ctx, cmName, d.ExecClient)
	return err
}

func (d *PoseidonSyncActivities) ResumeConsensusClient(ctx context.Context) error {
	cmName := fmt.Sprintf("cm-%s", d.ConsensusClient)
	_, err := d.BeaconActionsClient.StartClient(ctx, cmName, d.ConsensusClient)
	return err
}

// IsExecClientSynced only checks the first result
func (d *PoseidonSyncActivities) IsExecClientSynced(ctx context.Context) (bool, error) {
	syncStatuses, err := d.BeaconActionsClient.GetExecClientSyncStatus(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("SyncExecStatus")
		return false, err
	}
	for _, ss := range syncStatuses {
		log.Ctx(ctx).Info().Interface("syncStatus", ss)
		if ss.Result == false {
			return !ss.Result, nil
		}
	}
	return false, errors.New("not synced yet")
}

// IsConsensusClientSynced only checks the first result
func (d *PoseidonSyncActivities) IsConsensusClientSynced(ctx context.Context) (bool, error) {
	syncStatuses, err := d.BeaconActionsClient.GetConsensusClientSyncStatus(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("SyncExecStatus")
		return false, err
	}
	for _, ss := range syncStatuses {
		log.Ctx(ctx).Info().Interface("syncStatus", ss)
		if ss.Data.IsSyncing == false {
			return !ss.Data.IsSyncing, nil
		}
	}
	return false, errors.New("not synced yet")
}

func (d *PoseidonSyncActivities) RsyncExecBucket(ctx context.Context) error {
	br := poseidon_buckets.GethMainnetBucket
	err := d.Upload(ctx, br)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("RsyncExecBucket")
		return err
	}
	return err
}

func (d *PoseidonSyncActivities) RsyncConsensusBucket(ctx context.Context) error {
	br := poseidon_buckets.LighthouseMainnetBucket
	err := d.Upload(ctx, br)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("RsyncConsensusBucket")
		return err
	}
	return err
}
