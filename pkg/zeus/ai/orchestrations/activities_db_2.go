package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
)

func (z *ZeusAiPlatformActivities) SelectWorkflowIO(ctx context.Context, refID int) (WorkflowStageIO, error) {
	ws, err := artemis_orchestrations.SelectWorkflowStageReference(ctx, refID)
	if err != nil {
		log.Err(err).Interface("refID", refID).Msg("SelectWorkflowIO: failed to select workflow stage reference")
		return WorkflowStageIO{}, err
	}
	wsr := WorkflowStageIO{
		WorkflowStageReference: ws,
	}
	err = json.Unmarshal(ws.InputData, &wsr.WorkflowStageInfo)
	if err != nil {
		log.Err(err).Interface("ws", ws).Msg("SelectWorkflowIO: failed to unmarshal workflow stage info")
		return wsr, err
	}
	return wsr, nil
}

func (z *ZeusAiPlatformActivities) SaveWorkflowIO(ctx context.Context, wfInputs *WorkflowStageIO) (*WorkflowStageIO, error) {
	wsr := wfInputs.WorkflowStageReference
	b, err := json.Marshal(wfInputs.WorkflowStageInfo)
	if err != nil {
		log.Err(err).Interface("wfInputs", wfInputs).Msg("SaveWorkflowIO: failed to marshal workflow stage info")
		return nil, err
	}
	wsr.InputData = b
	err = artemis_orchestrations.InsertWorkflowStageReference(ctx, &wsr)
	if err != nil {
		log.Err(err).Interface("wfInputs", wfInputs).Msg("SaveWorkflowIO: failed to save workflow stage reference")
		return nil, err
	}
	wfInputs.WorkflowStageReference = wsr
	return wfInputs, nil
}

type WorkflowStageIO struct {
	artemis_orchestrations.WorkflowExecParams     `json:"workflowExecParams"`
	artemis_orchestrations.WorkflowStageReference `json:"workflowStageReference"`
	WorkflowStageInfo                             `json:"workflowStageInfo"`
	InputDataAnalysisToAgg                        `json:"inputDataAnalysisToAgg"`
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

func (ws *WorkflowStageInfo) GetOutSearchGroups() []*hera_search.SearchResultGroup {
	if ws.PromptReduction != nil && ws.PromptReduction.PromptReductionSearchResults != nil && len(ws.PromptReduction.PromptReductionSearchResults.OutSearchGroups) > 0 {
		return ws.PromptReduction.PromptReductionSearchResults.OutSearchGroups
	}
	return nil
}

func (ws *WorkflowStageInfo) GetOutTextGroups() []string {
	if ws.PromptReduction != nil && ws.PromptReduction.PromptReductionText != nil && len(ws.PromptReduction.PromptReductionText.OutPromptChunks) > 0 {
		return ws.PromptReduction.PromptReductionText.OutPromptChunks
	}
	return nil
}

func (ws *WorkflowStageIO) GetSearchGroupsOutByRetNameMatch(retNames map[string]bool) []hera_search.SearchResultGroup {
	if len(retNames) == 0 {
		return nil
	}
	if ws.PromptReduction == nil {
		return nil
	}
	if ws.PromptReduction.PromptReductionSearchResults == nil {
		return nil
	}
	var sgs []hera_search.SearchResultGroup
	for _, sgv := range ws.PromptReduction.PromptReductionSearchResults.OutSearchGroups {
		if sgv != nil && sgv.RetrievalName != nil {
			if tr, ok := retNames[*sgv.RetrievalName]; ok && tr {
				sgs = append(sgs, *sgv)
			}
		}
	}
	if len(sgs) <= 0 {
		log.Warn().Msg("GetSearchGroupsOutByRetNameMatch: sgs was empty")
	}
	return sgs
}
