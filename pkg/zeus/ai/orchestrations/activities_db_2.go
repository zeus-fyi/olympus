package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

func (z *ZeusAiPlatformActivities) SelectWorkflowIO(ctx context.Context, refID int) (WorkflowStageIO, error) {
	ws, err := artemis_orchestrations.SelectWorkflowStageReference(ctx, refID)
	if err != nil {
		log.Err(err).Interface("refID", refID).Msg("failed to select workflow stage reference")
		return WorkflowStageIO{}, err
	}
	wsr := WorkflowStageIO{
		WorkflowStageReference: ws,
	}
	err = json.Unmarshal(ws.InputData, &wsr.WorkflowStageInfo)
	if err != nil {
		log.Err(err).Interface("ws", ws).Msg("failed to unmarshal workflow stage info")
		return wsr, err
	}
	return wsr, nil
}

func (z *ZeusAiPlatformActivities) SaveWorkflowIO(ctx context.Context, wfInputs *WorkflowStageIO) error {
	wsr := wfInputs.WorkflowStageReference
	b, err := json.Marshal(wfInputs.WorkflowStageInfo)
	if err != nil {
		log.Err(err).Interface("wfInputs", wfInputs).Msg("failed to marshal workflow stage info")
		return err
	}
	wsr.InputData = b
	err = artemis_orchestrations.InsertWorkflowStageReference(ctx, &wsr)
	if err != nil {
		log.Err(err).Interface("wfInputs", wfInputs).Msg("failed to select workflow stage reference")
		return err
	}
	return nil
}

type WorkflowStageIO struct {
	artemis_orchestrations.WorkflowStageReference
	WorkflowStageInfo
}

type WorkflowStageInfo struct {
	TaskToExecute                      *TaskToExecute                      `json:"tte,omitempty"`
	RunAiWorkflowAutoEvalProcessInputs *RunAiWorkflowAutoEvalProcessInputs `json:"runAiWorkflowAutoEvalProcessInputs,omitempty"`
	CreateTriggerActionsWorkflowInputs *CreateTriggerActionsWorkflowInputs `json:"createTriggerActionsWorkflowInputs,omitempty"`
}
