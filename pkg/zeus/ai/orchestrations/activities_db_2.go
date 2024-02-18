package ai_platform_service_orchestrations

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (z *ZeusAiPlatformActivities) SelectWorkflowIO(ctx context.Context, ou org_users.OrgUser, wfID string) (WorkflowStageIO, error) {
	return WorkflowStageIO{}, nil
}

func (z *ZeusAiPlatformActivities) SaveWorkflowIO(ctx context.Context, ou org_users.OrgUser, wfInputs *WorkflowStageIO) (any, error) {
	return nil, nil
}

type WorkflowStageReference struct {
	InputID         int    `json:"inputID,omitempty"`
	InputStrID      string `json:"inputStrID,omitempty"`
	OrchestrationID int    `json:"orchestrationID,omitempty"`
	ChildWfID       string `json:"childWfID,omitempty"`
	RunCycle        int    `json:"runCycle"`
}

type WorkflowStageIO struct {
	WorkflowStageReference
	TaskToExecute                      *TaskToExecute                      `json:"tte,omitempty"`
	RunAiWorkflowAutoEvalProcessInputs *RunAiWorkflowAutoEvalProcessInputs `json:"runAiWorkflowAutoEvalProcessInputs,omitempty"`
	CreateTriggerActionsWorkflowInputs *CreateTriggerActionsWorkflowInputs `json:"createTriggerActionsWorkflowInputs,omitempty"`
}
