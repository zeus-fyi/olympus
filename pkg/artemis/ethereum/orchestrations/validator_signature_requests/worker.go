package eth_validator_signature_requests

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
	"go.temporal.io/sdk/client"
)

type ArtemisEthereumValidatorSignatureRequestsWorker struct {
	temporal_base.Worker
}

func (t *ArtemisEthereumValidatorSignatureRequestsWorker) ExecuteValidatorSignatureRequestsWorkflow(ctx context.Context, sigRequests aegis_inmemdbs.EthereumBLSKeySignatureRequests, signType string) (aegis_inmemdbs.EthereumBLSKeySignatureResponses, error) {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: fmt.Sprintf("%s-%s", t.TaskQueueName, signType),
	}
	sigReqWf := NewArtemisEthereumValidatorSignatureRequestWorkflow()
	wf := sigReqWf.ArtemisSendValidatorSignatureRequestsWorkflow
	workflowRun, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, sigRequests)
	if err != nil {
		log.Err(err).Msg("Hydra: Artemis Subsystem: ExecuteValidatorSignatureRequestsWorkflow")
		return aegis_inmemdbs.EthereumBLSKeySignatureResponses{}, err
	}
	var resp aegis_inmemdbs.EthereumBLSKeySignatureResponses
	err = workflowRun.Get(ctx, &resp)
	if err != nil {
		log.Err(err).Msg("Hydra: Artemis Subsystem: ExecuteValidatorSignatureRequestsWorkflow")
		return aegis_inmemdbs.EthereumBLSKeySignatureResponses{}, err
	}
	return resp, err
}
