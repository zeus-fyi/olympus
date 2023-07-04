package artemis_mev_transcations

import (
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	dynamodb_mev "github.com/zeus-fyi/olympus/datastores/dynamodb/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type HistoricalTxAnalysis struct {
	StartTimeDelay time.Duration
	Trades         []artemis_autogen_bases.EthMempoolMevTx
}

func (t *ArtemisMevWorkflow) ArtemisHistoricalSimTxWorkflow(ctx workflow.Context, trades HistoricalTxAnalysis) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 300,
	}
	srr := workflow.Sleep(ctx, trades.StartTimeDelay)
	if srr != nil {
		log.Error("Failed to sleep before tx analysis", "Error", srr)
		return srr
	}
	histSimTxCtx := workflow.WithActivityOptions(ctx, ao)
	for _, trade := range trades.Trades {
		err := workflow.ExecuteActivity(histSimTxCtx, t.HistoricalSimulateAndValidateTx, trade).Get(histSimTxCtx, nil)
		if err != nil {
			log.Error("Failed to sim historical mempool tx", "Error", err)
			return err
		}
	}
	return nil
}

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

func (t *ArtemisMevWorkflow) ArtemisRemoveProcessedTxsWorkflow(ctx workflow.Context, mempoolTxs []dynamodb_mev.MempoolTxsDynamoDB) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 12,
	}
	removeMempoolTxsCtx := workflow.WithActivityOptions(ctx, ao)
	for _, tx := range mempoolTxs {
		err := workflow.ExecuteActivity(removeMempoolTxsCtx, t.RemoveProcessedTx, tx).Get(removeMempoolTxsCtx, nil)
		if err != nil {
			log.Error("Failed to remove mempool txs", "Error", err)
			return err
		}
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
	var mempoolTxs []dynamodb_mev.MempoolTxsDynamoDB
	err := workflow.ExecuteActivity(getMempoolTxsCtx, t.GetMempoolTxs).Get(getMempoolTxsCtx, &mempoolTxs)
	if err != nil {
		log.Error("Failed to get mempool txs", "Error", err)
		return err
	}

	convertAo := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 60,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 20,
		},
	}
	var convertedMempoolTxs map[string]map[string]*types.Transaction
	convertMempoolTxsCtx := workflow.WithActivityOptions(ctx, convertAo)
	err = workflow.ExecuteActivity(convertMempoolTxsCtx, t.ConvertMempoolTxs, mempoolTxs).Get(convertMempoolTxsCtx, &convertedMempoolTxs)
	if err != nil {
		log.Error("Failed to convert mempool txs", "Error", err)
		return err
	}

	processAo := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 60,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 5,
		},
	}
	processMempoolTxsCtx := workflow.WithActivityOptions(ctx, processAo)
	var trades []artemis_autogen_bases.EthMempoolMevTx
	err = workflow.ExecuteActivity(processMempoolTxsCtx, t.ProcessMempoolTxs, convertedMempoolTxs).Get(processMempoolTxsCtx, &trades)
	if err != nil {
		log.Error("Failed to process mempool txs", "Error", err)
		return err
	}
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	childWorkflowFuture := workflow.ExecuteChildWorkflow(ctx, "ArtemisRemoveProcessedTxsWorkflow", mempoolTxs)
	var childWE workflow.Execution
	if err = childWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &childWE); err != nil {
		log.Error("Failed to get child workflow execution", "Error", err)
		return err
	}

	childWorkflowOptions = workflow.ChildWorkflowOptions{
		TaskQueue:         EthereumMainnetMevHistoricalTxTaskQueue,
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	histTxTrades := HistoricalTxAnalysis{
		StartTimeDelay: 12 * time.Second,
		Trades:         trades,
	}
	childWorkflowFutureHistoricalSimTx := workflow.ExecuteChildWorkflow(ctx, "ArtemisHistoricalSimTxWorkflow", histTxTrades)
	var childWEHistoricalSimTx workflow.Execution
	if err = childWorkflowFutureHistoricalSimTx.GetChildWorkflowExecution().Get(ctx, &childWEHistoricalSimTx); err != nil {
		log.Error("Failed to get sim historical tx workflow execution", "Error", err)
		return err
	}
	// Validate txs to bundle

	// Discard any bad txs

	// Simulate bundle

	// Create final bundle

	// Submit to flashbots before deadline
	return nil
}
