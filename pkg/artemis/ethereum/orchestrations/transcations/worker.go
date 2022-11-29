package artemis_ethereum_transcations

import (
	"context"

	"github.com/gochain/gochain/v4/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
	"go.temporal.io/sdk/client"
)

func (t *ArtemisEthereumTxWorker) ExecuteArtemisSendSignedTxWorkflow(ctx context.Context, params *types.Transaction) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewArtemisBroadcastEthereumTxWorkflow()
	wf := txWf.ArtemisSendSignedTxWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteArtemisSendSignedTxWorkflow")
		return err
	}
	return err
}

func (t *ArtemisEthereumTxWorker) ExecuteArtemisSendEthTxWorkflow(ctx context.Context, params web3_actions.SendTxPayload) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewArtemisBroadcastEthereumTxWorkflow()
	wf := txWf.ArtemisSendEthTxWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteArtemisSendEthTxWorkflow")
		return err
	}
	return err
}
