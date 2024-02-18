package ai_platform_service_orchestrations

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (z *ZeusAiPlatformActivities) SelectWorkflowIO(ctx context.Context, ou org_users.OrgUser, wfID string) (WorkflowStageIO, error) {
	return WorkflowStageIO{}, nil
}

func (z *ZeusAiPlatformActivities) SaveWorkflowIO(ctx context.Context, ou org_users.OrgUser, wfInputs *WorkflowStageIO) (any, error) {
	return nil, nil
}

type WorkflowStageIO struct {
	artemis_orchestrations.WorkflowStageReference
	TaskToExecute                      *TaskToExecute                      `json:"tte,omitempty"`
	RunAiWorkflowAutoEvalProcessInputs *RunAiWorkflowAutoEvalProcessInputs `json:"runAiWorkflowAutoEvalProcessInputs,omitempty"`
	CreateTriggerActionsWorkflowInputs *CreateTriggerActionsWorkflowInputs `json:"createTriggerActionsWorkflowInputs,omitempty"`
}
