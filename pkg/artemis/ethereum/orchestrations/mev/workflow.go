package artemis_mev_transcations

import (
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type ArtemisMevWorkflow struct {
	temporal_base.Workflow
	ArtemisMevActivities
}

const defaultTimeout = 6 * time.Second

func NewArtemisMevWorkflow() ArtemisMevWorkflow {
	deployWf := ArtemisMevWorkflow{
		Workflow:             temporal_base.Workflow{},
		ArtemisMevActivities: ArtemisMevActivities{},
	}
	return deployWf
}

func (t *ArtemisMevWorkflow) GetWorkflows() []interface{} {
	return []interface{}{t.ArtemisSendEthTxWorkflow, t.ArtemisSendSignedTxWorkflow, t.ArtemisMevWorkflow, t.ArtemisTxBlacklistWorkflow,
		t.ArtemisRemoveProcessedTxsWorkflow, t.ArtemisHistoricalSimTxWorkflow, t.ArtemisTxBlacklistProcessedTxsWorkflow,
		t.ArtemisGetLookaheadPricesWorkflow, t.GetTxReceipts,
	}
}

func (t *ArtemisMevWorkflow) ArtemisSendEthTxWorkflow(ctx workflow.Context, params web3_actions.SendEtherPayload) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	sendCtx := workflow.WithActivityOptions(ctx, ao)
	var txHash accounts.Hash
	err := workflow.ExecuteActivity(sendCtx, t.SendEther, params).Get(sendCtx, &txHash)
	if err != nil {
		log.Info("params", params)
		log.Info("txData", txHash)
		log.Error("Failed to send ether", "Error", err)
		return err
	}

	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 15,
		BackoffCoefficient: 2,
	}
	ao.RetryPolicy = retryPolicy
	rxCtx := workflow.WithActivityOptions(ctx, ao)
	var rx *types.Receipt

	err = workflow.ExecuteActivity(rxCtx, t.WaitForTxReceipt, txHash).Get(rxCtx, &rx)
	if err != nil {
		log.Info("params", params)
		log.Info("txData", txHash)
		log.Info("rx", rx)
		log.Error("Failed to get tx status", "Error", err)
		return err
	}
	return nil
}

func (t *ArtemisMevWorkflow) ArtemisSendSignedTxWorkflow(ctx workflow.Context, params *types.Transaction) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	sendCtx := workflow.WithActivityOptions(ctx, ao)
	var txData *types.Transaction
	err := workflow.ExecuteActivity(sendCtx, t.SubmitSignedTx, params).Get(sendCtx, &txData)
	if err != nil {
		log.Info("params", params)
		log.Info("txData", txData)
		log.Error("Failed to send signed tx", "Error", err)
		return err
	}
	rxCtx := workflow.WithActivityOptions(ctx, ao)
	var rx *types.Receipt
	err = workflow.ExecuteActivity(rxCtx, t.WaitForTxReceipt, txData.Hash()).Get(rxCtx, &rx)
	if err != nil {
		log.Info("params", params)
		log.Info("txData", txData)
		log.Info("rx", rx)
		log.Error("Failed to get tx status", "Error", err)
		return err
	}
	return nil
}
