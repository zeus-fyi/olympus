package poseidon_orchestrations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	pg_poseidon "github.com/zeus-fyi/olympus/datastores/postgres/apps/poseidon"
	athena_client "github.com/zeus-fyi/olympus/pkg/athena/client"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_buckets"
	beacon_actions "github.com/zeus-fyi/zeus/cookbooks/ethereum/beacons/actions"
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
	return []interface{}{d.PauseExecClient, d.PauseConsensusClient, d.ResumeExecClient, d.ResumeConsensusClient,
		d.IsExecClientSynced, d.IsConsensusClientSynced, d.RsyncExecBucket, d.RsyncConsensusBucket,
		d.ScheduleDiskWipe, d.RestartBeaconPod, d.ScheduleDiskUpload,
	}
}

func (d *PoseidonSyncActivities) PauseExecClient(ctx context.Context) error {
	cmName := fmt.Sprintf("cm-%s", PoseidonSyncActivitiesOrchestrator.ExecClient)
	_, err := PoseidonSyncActivitiesOrchestrator.BeaconActionsClient.PauseClient(ctx, cmName, PoseidonSyncActivitiesOrchestrator.ExecClient)
	return err
}

func (d *PoseidonSyncActivities) PauseConsensusClient(ctx context.Context) error {
	cmName := fmt.Sprintf("cm-%s", PoseidonSyncActivitiesOrchestrator.ConsensusClient)
	_, err := PoseidonSyncActivitiesOrchestrator.BeaconActionsClient.PauseClient(ctx, cmName, PoseidonSyncActivitiesOrchestrator.ConsensusClient)
	return err
}

func (d *PoseidonSyncActivities) ResumeExecClient(ctx context.Context) error {
	cmName := fmt.Sprintf("cm-%s", PoseidonSyncActivitiesOrchestrator.ExecClient)
	_, err := PoseidonSyncActivitiesOrchestrator.BeaconActionsClient.StartClient(ctx, cmName, PoseidonSyncActivitiesOrchestrator.ExecClient)
	return err
}

func (d *PoseidonSyncActivities) ResumeConsensusClient(ctx context.Context) error {
	cmName := fmt.Sprintf("cm-%s", PoseidonSyncActivitiesOrchestrator.ConsensusClient)
	_, err := PoseidonSyncActivitiesOrchestrator.BeaconActionsClient.StartClient(ctx, cmName, PoseidonSyncActivitiesOrchestrator.ConsensusClient)
	return err
}

// IsExecClientSynced only checks the first result
func (d *PoseidonSyncActivities) IsExecClientSynced(ctx context.Context, params pg_poseidon.UploadDataDirOrchestration) (bool, error) {
	bac := PoseidonSyncActivitiesOrchestrator.BeaconActionsClient
	bac.ExecClient = params.ClientName
	syncStatuses, err := bac.GetExecClientSyncStatus(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetExecClientSyncStatus")
		return false, err
	}
	if len(syncStatuses) <= 0 {
		return false, errors.New("no sync statuses returned")
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
func (d *PoseidonSyncActivities) IsConsensusClientSynced(ctx context.Context, params pg_poseidon.UploadDataDirOrchestration) (bool, error) {
	bac := PoseidonSyncActivitiesOrchestrator.BeaconActionsClient
	bac.ConsensusClient = params.ClientName
	syncStatuses, err := bac.GetConsensusClientSyncStatus(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetConsensusClientSyncStatus")
		return false, err
	}
	if len(syncStatuses) <= 0 {
		return false, errors.New("no sync statuses returned")
	}
	for _, ss := range syncStatuses {
		log.Ctx(ctx).Info().Interface("syncStatus", ss)
		if ss.Data.IsSyncing == false {
			return !ss.Data.IsSyncing, nil
		}
	}
	return false, errors.New("not synced yet")
}

type Response struct {
	Message string `json:"message"`
}

func (d *PoseidonSyncActivities) RsyncExecBucket(ctx context.Context) error {
	br := poseidon_buckets.GethMainnetBucket
	resp, err := PoseidonSyncActivitiesOrchestrator.UploadViaPortForward(ctx, PoseidonSyncActivitiesOrchestrator.BeaconKnsReq, br)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("RsyncExecBucket")
		return err
	}

	msg := Response{}
	if len(resp.ReplyBodies) <= 0 {
		return errors.New("not done")
	}
	for _, rep := range resp.ReplyBodies {
		err = json.Unmarshal(rep, &msg)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("GetConsensusClientSyncStatus")
			return err
		}
		if msg.Message != "done" {
			return errors.New("not done")
		}
	}
	return err
}

func (d *PoseidonSyncActivities) RsyncConsensusBucket(ctx context.Context) error {
	br := poseidon_buckets.LighthouseMainnetBucket
	_, err := PoseidonSyncActivitiesOrchestrator.UploadViaPortForward(ctx, PoseidonSyncActivitiesOrchestrator.BeaconKnsReq, br)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("RsyncConsensusBucket")
		return err
	}
	return err
}
