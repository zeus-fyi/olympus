package artemis_mev_transcations

import (
	"context"

	"github.com/ethereum/go-ethereum/v4/core/types"
	"github.com/ethereum/go-ethereum/web3/web3_actions"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func (t *ArtemisMevWorker) ExecuteArtemisMevWorkflow(ctx context.Context) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewArtemisMevWorkflow()
	wf := txWf.ArtemisMevWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf)
	if err != nil {
		log.Err(err).Msg("ExecuteArtemisMevWorkflow")
		return err
	}
	return err
}

func (t *ArtemisMevWorker) ExecuteArtemisSendSignedTxWorkflow(ctx context.Context, params *types.Transaction) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewArtemisMevWorkflow()
	wf := txWf.ArtemisSendSignedTxWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteArtemisSendSignedTxWorkflow")
		return err
	}
	return err
}

func (t *ArtemisMevWorker) ExecuteArtemisSendEthTxWorkflow(ctx context.Context, params web3_actions.SendEtherPayload) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewArtemisMevWorkflow()
	wf := txWf.ArtemisSendEthTxWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteArtemisSendEthTxWorkflow")
		return err
	}
	return err
}
