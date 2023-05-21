package artemis_mev_transcations

import (
	mempool_txs "github.com/zeus-fyi/olympus/datastores/dynamodb/mempool"
	"go.temporal.io/sdk/workflow"
)

// TODO update timeouts and finish this workflow, sync times with mainnet/goerli block times

func (t *ArtemisMevWorkflow) ArtemisMevWorkflow(ctx workflow.Context) error {

	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	getMempoolTxsCtx := workflow.WithActivityOptions(ctx, ao)
	var mempoolTxs []mempool_txs.MempoolTxsDynamoDB
	err := workflow.ExecuteActivity(getMempoolTxsCtx, t.GetMempoolTxs).Get(getMempoolTxsCtx, &mempoolTxs)
	if err != nil {
		log.Error("Failed to get mempool txs", "Error", err)
		return err
	}
	// Process mempool txs

	// Create tx bundle

	// Sim & validate tx bundle

	// Discard any bad txs

	// Submit to flashbots before deadline

	return nil
}
