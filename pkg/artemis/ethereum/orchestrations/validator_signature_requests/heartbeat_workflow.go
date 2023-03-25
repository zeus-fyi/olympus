package eth_validator_signature_requests

import (
	"context"
	"time"

	"github.com/oleiade/lane/v2"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	"go.temporal.io/sdk/workflow"
)

var HeartbeatQueue = lane.NewQueue[string]()

func (t *ArtemisEthereumValidatorSignatureRequestWorkflow) ValidatorsHeartbeatWorkflow(ctx workflow.Context, params interface{}) error {
	wfLog := workflow.GetLogger(ctx)
	localCtx := context.Background()
	serviceRoutes, err := artemis_validator_service_groups_models.SelectValidatorsServiceRoutes(localCtx)
	if err != nil {
		wfLog.Error("Failed to select validators to heartbeat", "error", err)
		return err
	}
	for _, vsrInfo := range serviceRoutes.GroupToServiceMap {
		wfLog.Info("Heartbeat to service", "GroupName", vsrInfo.GroupName)
		groupSize := len(serviceRoutes.GroupToPubKeySlice[vsrInfo.GroupName])
		wfLog.Info("Group key size", "GroupToPubKeySlice", groupSize)
		if groupSize < 100 {
			HeartbeatQueue.Enqueue(vsrInfo.GroupName)
		}
	}

	i := 30
	for {
		err = workflow.ExecuteActivity(ctx, t.SendHeartbeat).Get(ctx, nil)
		if err != nil {
			wfLog.Error("Failed to send heartbeat", "error", err)
		}
		err = workflow.Sleep(ctx, 30*time.Second)
		if err != nil {
			wfLog.Error("failed to sleep", "error", err)
		}
		i++
		if i >= 30 {
			wfLog.Info("Clearing heartbeat queue")
			i = 0
			for {
				ql := HeartbeatQueue.Size()
				if ql == 0 {
					break
				}
				_, qOk := HeartbeatQueue.Dequeue()
				if !qOk {
					continue
				}
			}
			serviceRoutes, err = artemis_validator_service_groups_models.SelectValidatorsServiceRoutes(localCtx)
			if err != nil {
				wfLog.Error("Failed to select validators to heartbeat", "error", err)
				return err
			}
			for _, vsrInfo := range serviceRoutes.GroupToServiceMap {
				wfLog.Info("Heartbeat to service", "GroupName", vsrInfo.GroupName)
				groupSize := len(serviceRoutes.GroupToPubKeySlice[vsrInfo.GroupName])
				wfLog.Info("Group key size", "GroupToPubKeySlice", groupSize)
				if groupSize < 100 {
					HeartbeatQueue.Enqueue(vsrInfo.GroupName)
				}
			}
		}
	}
}
