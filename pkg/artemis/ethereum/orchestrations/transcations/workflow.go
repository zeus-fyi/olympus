package artemis_ethereum_transcations

import (
	"time"

	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/gochain/v4/core/types"
	web3_types "github.com/zeus-fyi/gochain/web3/types"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/workflow"
)

type ArtemisEthereumTxBroadcastWorkflow struct {
	temporal_base.Workflow
	ArtemisEthereumBroadcastTxActivities
}

const defaultTimeout = 20 * time.Minute

func NewArtemisBroadcastEthereumTxWorkflow() ArtemisEthereumTxBroadcastWorkflow {
	deployWf := ArtemisEthereumTxBroadcastWorkflow{
		Workflow:                             temporal_base.Workflow{},
		ArtemisEthereumBroadcastTxActivities: ArtemisEthereumBroadcastTxActivities{},
	}
	return deployWf
}

func (t *ArtemisEthereumTxBroadcastWorkflow) GetWorkflows() interface{} {
	return t.ArtemisSendEthTxWorkflow
}

func (t *ArtemisEthereumTxBroadcastWorkflow) ArtemisSendEthTxWorkflow(ctx workflow.Context, params web3_actions.SendEtherPayload) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	sendCtx := workflow.WithActivityOptions(ctx, ao)
	var txHash common.Hash
	err := workflow.ExecuteActivity(sendCtx, t.SendEther, params).Get(sendCtx, &txHash)
	if err != nil {
		log.Info("params", params)
		log.Info("txData", txHash)
		log.Error("Failed to send ether", "Error", err)
		return err
	}

	rxCtx := workflow.WithActivityOptions(ctx, ao)
	var rx *web3_types.Receipt
	err = workflow.ExecuteActivity(rxCtx, t.WaitForReceipt, txHash).Get(rxCtx, &rx)
	if err != nil {
		log.Info("params", params)
		log.Info("txData", txHash)
		log.Info("rx", rx)
		log.Error("Failed to get tx status", "Error", err)
		return err
	}
	return nil
}

func (t *ArtemisEthereumTxBroadcastWorkflow) ArtemisSendSignedTxWorkflow(ctx workflow.Context, params *types.Transaction) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	sendCtx := workflow.WithActivityOptions(ctx, ao)
	var txData *web3_types.Transaction
	err := workflow.ExecuteActivity(sendCtx, t.SubmitSignedTxAndReturnTxData, params).Get(sendCtx, &txData)
	if err != nil {
		log.Info("params", params)
		log.Info("txData", txData)
		log.Error("Failed to send signed tx", "Error", err)
		return err
	}
	rxCtx := workflow.WithActivityOptions(ctx, ao)
	var rx *web3_types.Receipt
	err = workflow.ExecuteActivity(rxCtx, t.WaitForReceipt, txData.Hash).Get(rxCtx, &rx)
	if err != nil {
		log.Info("params", params)
		log.Info("txData", txData)
		log.Info("rx", rx)
		log.Error("Failed to get tx status", "Error", err)
		return err
	}
	return nil
}
