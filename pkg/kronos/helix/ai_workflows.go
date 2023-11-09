package kronos_helix

import (
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (k *KronosWorkflow) AiWorkflow(ctx workflow.Context, ou org_users.OrgUser, email, content string) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
		},
	}
	runAiTaskCtx := workflow.WithActivityOptions(ctx, ao)
	var resp openai.ChatCompletionResponse
	err := workflow.ExecuteActivity(runAiTaskCtx, k.AiTask, ou, content).Get(runAiTaskCtx, &resp)
	if err != nil {
		logger.Error("failed to execute AiTask", "Error", err)
		// You can decide if you want to return the error or continue monitoring.
		return err
	}

	sendEmailTaskCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(sendEmailTaskCtx, k.SendTaskResponseEmail, email, resp).Get(sendEmailTaskCtx, &resp)
	if err != nil {
		logger.Error("failed to execute SaveAiTaskResponse", "Error", err)
		// You can decide if you want to return the error or continue monitoring.
		return err
	}

	if ou.OrgID > 0 && ou.UserID > 0 {
		saveAiTaskCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(saveAiTaskCtx, k.SaveAiTaskResponse, ou, resp).Get(saveAiTaskCtx, &resp)
		if err != nil {
			logger.Error("failed to execute SaveAiTaskResponse", "Error", err)
			// You can decide if you want to return the error or continue monitoring.
			return err
		}
	}
	return nil
}
