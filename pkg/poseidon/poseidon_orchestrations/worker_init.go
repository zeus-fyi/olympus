package poseidon_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	beacon_actions "github.com/zeus-fyi/olympus/cookbooks/ethereum/beacons/actions"
	athena_client "github.com/zeus-fyi/olympus/pkg/athena/client"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

type PoseidonWorker struct {
	temporal_base.Worker
	beacon_actions.BeaconActionsClient
	athena_client.AthenaClient
}

var PoseidonSyncWorker PoseidonWorker
var PoseidonBearer string

const PoseidonTaskQueue = "PoseidonTaskQueue"

func InitPoseidonWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Poseidon: InitPoseidonWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("Poseidon: sync failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := PoseidonTaskQueue

	ba := beacon_actions.NewDefaultBeaconActionsClient(PoseidonBearer)

	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewPoseidonSyncActivity(ba)
	wf := NewPoseidonSyncWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	PoseidonSyncWorker.Worker = w
	PoseidonSyncWorker.TemporalClient = tc
	return
}
