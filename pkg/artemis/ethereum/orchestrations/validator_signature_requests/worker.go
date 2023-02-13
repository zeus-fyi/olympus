package eth_validator_signature_requests

import (
	"context"
	"github.com/rs/zerolog/log"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
	bls_signer "github.com/zeus-fyi/zeus/pkg/crypto/bls"
	"go.temporal.io/sdk/client"
)

type ArtemisEthereumValidatorSignatureRequestsWorker struct {
	temporal_base.Worker
}

func init() {
	_ = bls_signer.InitEthBLS()
}

func (t *ArtemisEthereumValidatorSignatureRequestsWorker) ExecuteValidatorSignatureRequestsWorkflow(ctx context.Context, sigRequests aegis_inmemdbs.EthereumBLSKeySignatureRequests) (aegis_inmemdbs.EthereumBLSKeySignatureResponses, error) {
	if len(sigRequests.Map) == 0 {
		return aegis_inmemdbs.EthereumBLSKeySignatureResponses{}, nil
	}
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
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
