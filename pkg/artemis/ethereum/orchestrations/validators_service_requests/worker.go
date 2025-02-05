package eth_validators_service_requests

import (
	"context"

	bls_signer "github.com/zeus-fyi/zeus/pkg/aegis/crypto/bls"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"

	"github.com/rs/zerolog/log"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"go.temporal.io/sdk/client"
)

type ArtemisEthereumValidatorsRequestsWorker struct {
	temporal_base.Worker
}

const (
	EthereumMainnetValidatorsRequestsTaskQueue  = "EthereumMainnetValidatorsRequestsTaskQueue"
	EthereumGoerliValidatorsRequestsTaskQueue   = "EthereumGoerliValidatorsRequestsTaskQueue"
	EthereumEphemeryValidatorsRequestsTaskQueue = "EthereumEphemeryValidatorsRequestsTaskQueue"
)

func init() {
	_ = bls_signer.InitEthBLS()
}

var (
	ArtemisEthereumMainnetValidatorsRequestsWorker ArtemisEthereumValidatorsRequestsWorker
	MainnetStakingCloudCtxNs                       = zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "mainnet-staking",
		Env:           "production",
	}
	ArtemisEthereumGoerliValidatorsRequestsWorker ArtemisEthereumValidatorsRequestsWorker
	GoerliStakingCloudCtxNs                       = zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "goerli-staking",
		Env:           "production",
	}
	ArtemisEthereumEphemeryValidatorsRequestsWorker ArtemisEthereumValidatorsRequestsWorker
	EphemeryStakingCloudCtxNs                       = zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "ephemeral-staking",
		Env:           "production",
	}
)

type ValidatorServiceGroupWorkflowRequest struct {
	OrgID int
	hestia_req_types.ServiceRequestWrapper
	hestia_req_types.ValidatorServiceOrgGroupSlice
}

func (t *ArtemisEthereumValidatorsRequestsWorker) ExecuteServiceNewValidatorsToCloudCtxNsWorkflow(ctx context.Context, params ValidatorServiceGroupWorkflowRequest) error {
	log.Info().Msg("ArtemisEthereumValidatorsRequestsWorker: ExecuteServiceNewValidatorsToCloudCtxNsWorkflow")
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	vsWf := NewArtemisEthereumValidatorServiceRequestWorkflow()
	wf := vsWf.ServiceNewValidatorsToCloudCtxNsWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteServiceNewValidatorsToCloudCtxNsWorkflow: ServiceNewValidatorsToCloudCtxNsWorkflow")
		return err
	}
	return err
}
