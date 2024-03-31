package ai_platform_service_orchestrations

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

type WorkflowStageIO struct {
	artemis_orchestrations.WorkflowExecParams     `json:"workflowExecParams"`
	artemis_orchestrations.WorkflowStageReference `json:"workflowStageReference"`
	WorkflowStageInfo                             `json:"workflowStageInfo"`
}

type WorkflowStageInfo struct {
	ApiIterationCount                  int                                 `json:"apiIterationCount"`
	Metadata                           json.RawMessage                     `json:"metadata,omitempty"`
	WorkflowInCacheHash                map[string]bool                     `json:"workflowInCacheHash,omitempty"`
	RunAiWorkflowAutoEvalProcessInputs *RunAiWorkflowAutoEvalProcessInputs `json:"runAiWorkflowAutoEvalProcessInputs,omitempty"`
	CreateTriggerActionsWorkflowInputs *CreateTriggerActionsWorkflowInputs `json:"createTriggerActionsWorkflowInputs,omitempty"`
	PromptReduction                    *PromptReduction                    `json:"promptReduction,omitempty"`
	PromptTextFromTextStage            string                              `json:"promptTextFromTextStage,omitempty"`
}
