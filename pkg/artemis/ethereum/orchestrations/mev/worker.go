package artemis_mev_transcations

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"go.temporal.io/sdk/client"
)

func (t *ArtemisMevWorker) ExecuteArtemisMevWorkflow(ctx context.Context) error {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewArtemisMevWorkflow()
	wf := txWf.ArtemisMevWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf)
	if err != nil {
		log.Err(err).Msg("ExecuteArtemisMevWorkflow")
		return err
	}
	return err
}

func (t *ArtemisMevWorker) ExecuteArtemisBlacklistTxWorkflow(ctx context.Context) error {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewArtemisMevWorkflow()
	wf := txWf.ArtemisTxBlacklistWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf)
	if err != nil {
		log.Err(err).Msg("ExecuteArtemisBlacklistTxWorkflow")
		return err
	}
	return err
}

func (t *ArtemisMevWorker) ExecuteArtemisSendSignedTxWorkflow(ctx context.Context, params *types.Transaction) error {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewArtemisMevWorkflow()
	wf := txWf.ArtemisSendSignedTxWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteArtemisSendSignedTxWorkflow")
		return err
	}
	return err
}

func (t *ArtemisMevWorker) ExecuteArtemisSendEthTxWorkflow(ctx context.Context, params web3_actions.SendEtherPayload) error {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewArtemisMevWorkflow()
	wf := txWf.ArtemisSendEthTxWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteArtemisSendEthTxWorkflow")
		return err
	}
	return err
}
