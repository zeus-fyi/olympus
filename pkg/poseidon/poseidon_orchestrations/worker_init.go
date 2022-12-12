package poseidon_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	beacon_actions "github.com/zeus-fyi/olympus/cookbooks/ethereum/beacons/actions"
	client_consts "github.com/zeus-fyi/olympus/cookbooks/ethereum/beacons/constants"
	athena_client "github.com/zeus-fyi/olympus/pkg/athena/client"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
)

type PoseidonWorker struct {
	temporal_base.Worker
}

var PoseidonSyncWorker PoseidonWorker
var PoseidonBearer string

const PoseidonTaskQueue = "PoseidonTaskQueue"

var kCtxNsHeader = zeus_req_types.TopologyDeployRequest{
	TopologyID: 1669159384971627008,
	CloudCtxNs: zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "ethereum",
		Env:           "production",
	},
}

func InitPoseidonWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Poseidon: InitPoseidonWorker starting")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("Poseidon: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := PoseidonTaskQueue

	ba := beacon_actions.NewDefaultBeaconActionsClient(PoseidonBearer, kCtxNsHeader)
	ac := athena_client.NewLocalAthenaClient(PoseidonBearer)
	w := temporal_base.NewWorker(taskQueueName)

	PoseidonSyncActivitiesOrchestrator = NewPoseidonSyncActivity(ba, ac)
	PoseidonSyncActivitiesOrchestrator.BeaconActionsClient.ExecClient = client_consts.Geth
	PoseidonSyncActivitiesOrchestrator.BeaconActionsClient.ConsensusClient = client_consts.Lighthouse

	wf := NewPoseidonSyncWorkflow(PoseidonSyncActivitiesOrchestrator)

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(PoseidonSyncActivitiesOrchestrator.GetActivities())
	PoseidonSyncWorker.Worker = w
	PoseidonSyncWorker.TemporalClient = tc
	log.Ctx(ctx).Info().Msg("Poseidon: InitPoseidonWorker finished")
	return
}
