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

func (t *ArtemisEthereumValidatorSignatureRequestsWorker) ExecuteValidatorSignatureRequestsWorkflow(ctx context.Context, params any) (aegis_inmemdbs.EthereumBLSKeySignatureResponses, error) {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	sigReqWf := NewArtemisEthereumValidatorSignatureRequestActivities()
	wf, err := sigReqWf.RequestValidatorSignature(ctx, params)
	resp, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteArtemisSendSignedTxWorkflow")
		return aegis_inmemdbs.EthereumBLSKeySignatureResponses{}, err
	}
	fmt.Print(resp)
	return aegis_inmemdbs.EthereumBLSKeySignatureResponses{}, err
}
