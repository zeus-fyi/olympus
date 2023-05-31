package artemis_mev_transcations

import (
	"github.com/ethereum/go-ethereum/core/types"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"go.temporal.io/sdk/workflow"
)

func (t *ArtemisMevWorkflow) ArtemisTxBlacklistWorkflow(ctx workflow.Context) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	getMempoolTxsCtx := workflow.WithActivityOptions(ctx, ao)
	var mempoolTxs map[string]map[string]*types.Transaction
	err := workflow.ExecuteActivity(getMempoolTxsCtx, t.BlacklistMinedTxs).Get(getMempoolTxsCtx, &mempoolTxs)
	if err != nil {
		log.Error("Failed to get mempool txs", "Error", err)
		return err
	}
	return nil
}

// TODO update timeouts and finish this workflow, sync times with mainnet/goerli block times

func (t *ArtemisMevWorkflow) ArtemisMevWorkflow(ctx workflow.Context) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	getMempoolTxsCtx := workflow.WithActivityOptions(ctx, ao)
	var mempoolTxs map[string]map[string]*types.Transaction
	err := workflow.ExecuteActivity(getMempoolTxsCtx, t.GetMempoolTxs).Get(getMempoolTxsCtx, &mempoolTxs)
	if err != nil {
		log.Error("Failed to get mempool txs", "Error", err)
		return err
	}
	processMempoolTxsCtx := workflow.WithActivityOptions(ctx, ao)
	var trades []artemis_autogen_bases.EthMempoolMevTx
	err = workflow.ExecuteActivity(processMempoolTxsCtx, t.ProcessMempoolTxs, mempoolTxs).Get(processMempoolTxsCtx, &trades)
	if err != nil {
		log.Error("Failed to process mempool txs", "Error", err)
		return err
	}
	// Validate txs to bundle

	// Discard any bad txs

	// Simulate bundle

	// Create final bundle

	// Submit to flashbots before deadline
	return nil
}
