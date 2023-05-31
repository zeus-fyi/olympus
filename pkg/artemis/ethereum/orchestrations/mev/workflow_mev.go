package artemis_mev_transcations

import (
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	mempool_txs "github.com/zeus-fyi/olympus/datastores/dynamodb/mempool"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"go.temporal.io/sdk/workflow"
)

func (t *ArtemisMevWorkflow) ArtemisTxBlacklistWorkflow(ctx workflow.Context) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 12,
	}
	blacklistTxsCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(blacklistTxsCtx, t.BlacklistMinedTxs).Get(blacklistTxsCtx, nil)
	if err != nil {
		log.Error("Failed to blacklist mempool txs", "Error", err)
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
	var mempoolTxs []mempool_txs.MempoolTxsDynamoDB
	err := workflow.ExecuteActivity(getMempoolTxsCtx, t.GetMempoolTxs).Get(getMempoolTxsCtx, &mempoolTxs)
	if err != nil {
		log.Error("Failed to get mempool txs", "Error", err)
		return err
	}

	var convertedMempoolTxs map[string]map[string]*types.Transaction
	convertMempoolTxsCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(convertMempoolTxsCtx, t.ConvertMempoolTxs, mempoolTxs).Get(convertMempoolTxsCtx, &convertedMempoolTxs)
	if err != nil {
		log.Error("Failed to convert mempool txs", "Error", err)
		return err
	}

	processMempoolTxsCtx := workflow.WithActivityOptions(ctx, ao)
	var trades []artemis_autogen_bases.EthMempoolMevTx
	err = workflow.ExecuteActivity(processMempoolTxsCtx, t.ProcessMempoolTxs, convertedMempoolTxs).Get(processMempoolTxsCtx, &trades)
	if err != nil {
		log.Error("Failed to process mempool txs", "Error", err)
		return err
	}

	removeMempoolTxsCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(removeMempoolTxsCtx, t.RemoveProcessedTxs, mempoolTxs).Get(removeMempoolTxsCtx, nil)
	if err != nil {
		log.Error("Failed to remove mempool txs", "Error", err)
		return err
	}

	// Validate txs to bundle

	// Discard any bad txs

	// Simulate bundle

	// Create final bundle

	// Submit to flashbots before deadline
	return nil
}
