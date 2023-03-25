package eth_validator_signature_requests

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"go.temporal.io/sdk/client"
)

const HeartbeatTaskQueue = "HeartbeatTaskQueue"

var ArtemisEthereumValidatorSignatureRequestsHeartbeatWorker ArtemisEthereumValidatorSignatureRequestsWorker

func InitHeartbeatWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitHeartbeatWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitHeartbeatWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := HeartbeatTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisEthereumValidatorSignatureRequestActivities()
	wf := NewArtemisEthereumValidatorSignatureRequestWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisEthereumValidatorSignatureRequestsHeartbeatWorker.Worker = w
	ArtemisEthereumValidatorSignatureRequestsHeartbeatWorker.TemporalClient = tc
	return
}

func (t *ArtemisEthereumValidatorSignatureRequestsWorker) ExecuteHeartbeatWorkflow(ctx context.Context) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowID := "heartbeat"

	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: t.TaskQueueName,
	}
	sigReqWf := NewArtemisEthereumValidatorSignatureRequestWorkflow()
	wf := sigReqWf.ValidatorsHeartbeatWorkflow
	_, err := c.SignalWithStartWorkflow(ctx, workflowID, "start", nil, workflowOptions, wf, nil)
	if err != nil {
		log.Err(err).Msg("Hydra: Artemis Subsystem: ExecuteHeartbeatWorkflow")
		return err
	}
	return err
}
