package orchestrations

import (
	"context"

	beacon_actions "github.com/zeus-fyi/olympus/cookbooks/ethereum/beacons/actions"
	athena_client "github.com/zeus-fyi/olympus/pkg/athena/client"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_buckets"
)

type PoseidonSyncActivities struct {
	beacon_actions.BeaconActionsClient
}

func NewPoseidonSyncActivity(client beacon_actions.BeaconActionsClient) PoseidonSyncActivities {
	return PoseidonSyncActivities{client}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *PoseidonSyncActivities) GetActivities() ActivitiesSlice {
	return []interface{}{d.Pause, d.Resume, d.SyncExecStatus, d.SyncConsensusStatus, d.SyncStatus, d.RsyncExecBucket, d.RsyncConsensusBucket}
}

func (d *PoseidonSyncActivities) Pause(ctx context.Context, cmName, clientName string) error {
	_, err := d.BeaconActionsClient.PauseClient(ctx, cmName, clientName)
	return err
}

func (d *PoseidonSyncActivities) Resume(ctx context.Context, cmName, clientName string) error {
	_, err := d.BeaconActionsClient.StartClient(ctx, cmName, clientName)
	return err
}

// TODO convert these from bytes to struct values from json
func (d *PoseidonSyncActivities) SyncExecStatus(ctx context.Context) error {
	_, err := d.BeaconActionsClient.GetExecClientSyncStatus(ctx)
	return err
}

func (d *PoseidonSyncActivities) SyncConsensusStatus(ctx context.Context) error {
	_, err := d.BeaconActionsClient.GetConsensusClientSyncStatus(ctx)
	return err
}

func (d *PoseidonSyncActivities) SyncStatus(ctx context.Context, clientName string) error {
	return nil
}

func (d *PoseidonSyncActivities) RsyncExecBucket(ctx context.Context) error {
	ac := athena_client.NewLocalAthenaClient(PoseidonBearer)
	br := poseidon_buckets.GethMainnetBucket
	err := ac.Upload(ctx, br)
	return err
}

func (d *PoseidonSyncActivities) RsyncConsensusBucket(ctx context.Context) error {
	ac := athena_client.NewLocalAthenaClient(PoseidonBearer)
	br := poseidon_buckets.LighthouseMainnetBucket
	err := ac.Upload(ctx, br)
	return err
}
