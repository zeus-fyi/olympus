package artemis_mev_transcations

import (
	web3_types "github.com/ethereum/go-ethereum/web3/types"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"go.temporal.io/sdk/workflow"
)

// TODO update timeouts and finish this workflow, sync times with mainnet/goerli block times

func (t *ArtemisMevWorkflow) ArtemisMevWorkflow(ctx workflow.Context) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	getMempoolTxsCtx := workflow.WithActivityOptions(ctx, ao)
	var mempoolTxs map[string]map[string]*web3_types.RpcTransaction
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
