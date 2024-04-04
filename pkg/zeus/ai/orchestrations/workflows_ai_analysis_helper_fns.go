package ai_platform_service_orchestrations

import (
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func getDefaultAnalysisRetryPolicy() workflow.ActivityOptions {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Hour * 24, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 3,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    100,
		},
	}
	return ao
}

func getRetrievalWfRetryPolicy() workflow.ActivityOptions {
	aoRet := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Hour * 24, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 3,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    10000,
		},
	}
	return aoRet
}

func getAnalysisPrInput(analysisInst artemis_orchestrations.WorkflowTemplateData) *PromptReduction {
	pr := &PromptReduction{
		MarginBuffer:          analysisInst.AnalysisMarginBuffer,
		Model:                 analysisInst.AnalysisModel,
		TokenOverflowStrategy: analysisInst.AnalysisTokenOverflowStrategy,
		PromptReductionText: &PromptReductionText{
			InPromptBody: analysisInst.AnalysisPrompt,
		},
	}
	return pr
}

func getDummyChatCompResp() *ChatCompletionQueryResponse {
	aiResp := &ChatCompletionQueryResponse{
		Prompt: make(map[string]string),
		Response: openai.ChatCompletionResponse{
			Model: "none",
			Usage: openai.Usage{
				PromptTokens:     0,
				CompletionTokens: 0,
				TotalTokens:      0,
			},
		},
	}
	return aiResp
}

func getAnalysisTaskContext(analysisInst artemis_orchestrations.WorkflowTemplateData) TaskContext {
	return TaskContext{
		TaskName:              analysisInst.AnalysisTaskName,
		TaskType:              AnalysisTask,
		Model:                 analysisInst.AnalysisModel,
		TaskID:                analysisInst.AnalysisTaskID,
		ResponseFormat:        analysisInst.AnalysisResponseFormat,
		Prompt:                analysisInst.AnalysisPrompt,
		WorkflowTemplateData:  analysisInst,
		TokenOverflowStrategy: analysisInst.AnalysisTokenOverflowStrategy,
		MarginBuffer:          analysisInst.AnalysisMarginBuffer,
		Temperature:           float32(analysisInst.AnalysisTemperature),
	}
}
