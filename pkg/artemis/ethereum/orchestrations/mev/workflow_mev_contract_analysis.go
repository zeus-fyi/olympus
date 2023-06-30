package artemis_mev_transcations

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func (t *ArtemisMevWorkflow) ArtemisERC20TokenInfoFetchWorkflow(ctx workflow.Context, params any) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 300,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(ctx, t.FetchERC20TokenInfo, params).Get(ctx, nil)
	if err != nil {
		log.Error("Failed to ", "Error", err)
		return err
	}
	err = workflow.ExecuteActivity(ctx, t.FetchERC20TokenBalanceOfStorageSlot, params).Get(ctx, nil)
	if err != nil {
		log.Error("Failed to ", "Error", err)
		return err
	}
	err = workflow.ExecuteActivity(ctx, t.CalculateTransferTaxFee, params).Get(ctx, nil)
	if err != nil {
		log.Error("Failed to ", "Error", err)
		return err
	}
	return nil
}
