package artemis_ethereum_transcations

import (
	"context"

	"github.com/ethereum/go-ethereum/v4/core/types"
	"github.com/ethereum/go-ethereum/web3/web3_actions"
	"github.com/rs/zerolog/log"
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

func (t *ArtemisEthereumTxWorker) ExecuteArtemisSendEthTxWorkflow(ctx context.Context, params web3_actions.SendEtherPayload) error {
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
