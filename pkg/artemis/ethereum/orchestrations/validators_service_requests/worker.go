package eth_validators_service_requests

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"go.temporal.io/sdk/client"
)

type ArtemisEthereumValidatorsRequestsWorker struct {
	temporal_base.Worker
}

var (
	ArtemisEthereumMainnetValidatorsRequestsWorker  ArtemisEthereumValidatorsRequestsWorker
	ArtemisEthereumEphemeryValidatorsRequestsWorker ArtemisEthereumValidatorsRequestsWorker
)

const (
	EthereumMainnetValidatorsRequestsTaskQueue  = "EthereumMainnetValidatorsRequestsTaskQueue"
	EthereumEphemeryValidatorsRequestsTaskQueue = "EthereumEphemeryValidatorsRequestsTaskQueue"
)

type ValidatorServiceGroupWorkflowRequest struct {
	hestia_req_types.ServiceRequestWrapper
	hestia_req_types.ValidatorServiceOrgGroupSlice
}

func (t *ArtemisEthereumValidatorsRequestsWorker) ExecuteServiceNewValidatorsToCloudCtxNsWorkflow(ctx context.Context, params ValidatorServiceGroupWorkflowRequest) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewArtemisEthereumValidatorServiceRequestWorkflow()
	wf := txWf.ServiceNewValidatorsToCloudCtxNsWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ServiceNewValidatorsToCloudCtxNsWorkflow")
		return err
	}
	return err
}
