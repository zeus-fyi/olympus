package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/openai"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (z *ZeusAiPlatformActivities) CreateWsr(ctx context.Context, cp *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	// todo if exists already, skip
	wio := WorkflowStageIO{
		WorkflowStageReference: cp.Wsr,
		WorkflowExecParams:     cp.WfExecParams,
		WorkflowStageInfo: WorkflowStageInfo{
			PromptReduction: &PromptReduction{
				MarginBuffer:          cp.Tc.MarginBuffer,
				Model:                 cp.Tc.Model,
				TokenOverflowStrategy: cp.Tc.TokenOverflowStrategy,
			},
		},
	}
	wio.Org.OrgID = cp.Ou.OrgID
	_, err := s3ws(ctx, cp, &wio)
	if err != nil {
		log.Err(err).Msg("CreateWsr: failed")
		return nil, err
	}
	return cp, nil
}

func (z *ZeusAiPlatformActivities) SaveTaskOutput(ctx context.Context, wr *artemis_orchestrations.AIWorkflowAnalysisResult, cp *MbChildSubProcessParams, dataIn InputDataAnalysisToAgg) (int, error) {
	if cp == nil {
		return 0, fmt.Errorf("SaveTaskOutput: cp is nil")
	}
	if wr == nil {
		return 0, nil
	}
	wio, werr := gs3wfs(ctx, cp)
	if werr != nil {
		log.Err(werr).Msg("TokenOverflowReduction: failed to select workflow io")
		return 0, werr
	}
	osg := wio.GetOutSearchGroups()
	if osg != nil {
		if cp.Wsr.ChunkOffset < len(osg) {
			dataIn.SearchResultGroup = wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups[cp.Wsr.ChunkOffset]
		}
	}
	err := artemis_orchestrations.InsertAiWorkflowAnalysisResult(ctx, wr)
	if err != nil {
		log.Err(err).Interface("wr", wr).Interface("wr", wr).Msg("SaveTaskOutput: failed")
		return 0, err
	}
	err = s3wsCustomTaskName(ctx, cp, fmt.Sprintf("%d", wr.WorkflowResultID), dataIn)
	if err != nil {
		log.Err(err).Msg("s3wsCustomTaskName: failed")
		return -1, err
	}
	if dataIn.TextInput != nil {
		wio.PromptTextFromTextStage += *dataIn.TextInput
	}
	_, eerr := s3ws(ctx, cp, wio)
	if eerr != nil {
		log.Err(err).Msg("SaveTaskOutput: failed")
		return -1, eerr
	}
	return wr.WorkflowResultID, nil
}

// UpdateTaskOutput updates the task output, but it only intended for json output results
func (z *ZeusAiPlatformActivities) UpdateTaskOutput(ctx context.Context, cp *MbChildSubProcessParams) ([]artemis_orchestrations.JsonSchemaDefinition, error) {
	if cp == nil || len(cp.Tc.JsonResponseResults) <= 0 {
		return nil, nil
	}
	var skipAnalysis bool
	jro := FilterPassingEvalPassingResponses(cp.Tc.JsonResponseResults)
	var md []byte
	var err error
	var filteredJsonResponses []artemis_orchestrations.JsonSchemaDefinition
	var infoJsonResponses []artemis_orchestrations.JsonSchemaDefinition
	for evalState, v := range jro {
		switch evalState {
		case filterState:
			filteredJsonResponses = v.Passed
			log.Info().Interface("v.Passed", len(v.Passed)).Interface("v.Failed", len(v.Failed)).Msg("UpdateTaskOutput: filterState")
		case infoState:
			if len(v.Failed) > 0 {
				skipAnalysis = true
			} else {
				infoJsonResponses = v.Passed
				log.Info().Interface("v.Passed", len(v.Passed)).Msg("UpdateTaskOutput: infoState")
			}
		case errorState:
			// TODO: stop workflow?
		}
	}
	var sg *hera_search.SearchResultGroup
	wio, werr := gs3wfs(ctx, cp)
	if werr != nil {
		log.Err(werr).Msg("TokenOverflowReduction: failed to select workflow io")
		return nil, werr
	}
	if wio.PromptReduction != nil && wio.PromptReduction.PromptReductionSearchResults != nil && wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups != nil && len(wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups) > 0 {
		if cp.Wsr.ChunkOffset < len(wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups) {
			sg = wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups[cp.Wsr.ChunkOffset]
		}
	}
	var res []artemis_orchestrations.JsonSchemaDefinition
	if len(filteredJsonResponses) <= 0 && len(infoJsonResponses) <= 0 {
		skipAnalysis = true
		md, err = json.Marshal(jro)
		if err != nil {
			log.Err(err).Interface("jro", jro).Msg("UpdateTaskOutput: failed")
			return nil, err
		}
	} else if len(filteredJsonResponses) > 0 {
		res = filteredJsonResponses
		tmp := InputDataAnalysisToAgg{
			SearchResultGroup: sg,
			ChatCompletionQueryResponse: &ChatCompletionQueryResponse{
				JsonResponseResults: res,
				RegexSearchResults:  cp.Tc.RegexSearchResults,
			},
		}
		md, err = json.Marshal(tmp)
		if err != nil {
			log.Err(err).Interface("infoJsonResponses", infoJsonResponses).Interface("jro", jro).Msg("UpdateTaskOutput: failed")
			return nil, err
		}
	} else {
		res = infoJsonResponses
		tmp := InputDataAnalysisToAgg{
			SearchResultGroup: sg,
			ChatCompletionQueryResponse: &ChatCompletionQueryResponse{
				JsonResponseResults: res,
				RegexSearchResults:  cp.Tc.RegexSearchResults,
			},
		}
		md, err = json.Marshal(tmp)
		if err != nil {
			log.Err(err).Interface("infoJsonResponses", infoJsonResponses).Interface("jro", jro).Msg("UpdateTaskOutput: failed")
			return nil, err
		}
	}
	if res != nil && sg != nil && sg.SearchResults != nil {
		seen := make(map[int]bool)
		for _, jr := range res {
			for _, fv := range jr.Fields {
				if fv.FieldName == "msg_id" && fv.IsValidated && fv.NumberValue != nil && *fv.NumberValue > 0 {
					seen[int(*fv.NumberValue)] = true
				}
				if fv.FieldName == "msg_id" && fv.IsValidated && fv.IntegerValue != nil && *fv.IntegerValue > 0 {
					seen[*fv.IntegerValue] = true
				}
			}
		}
		sg.FilteredSearchResults = []hera_search.SearchResult{}
		for _, sr := range sg.SearchResults {
			_, ok := seen[sr.UnixTimestamp]
			if ok {
				sg.FilteredSearchResults = append(sg.FilteredSearchResults, sr)
			}
		}
	}
	wr := &artemis_orchestrations.AIWorkflowAnalysisResult{
		WorkflowResultID:      cp.Tc.WorkflowResultID,
		ResponseID:            cp.Tc.ResponseID,
		OrchestrationID:       cp.Oj.OrchestrationID,
		SourceTaskID:          cp.Tc.TaskID,
		IterationCount:        cp.Wsr.IterationCount,
		ChunkOffset:           cp.Wsr.ChunkOffset,
		RunningCycleNumber:    cp.Wsr.RunCycle,
		SearchWindowUnixStart: cp.Window.UnixStartTime,
		SearchWindowUnixEnd:   cp.Window.UnixEndTime,
		Metadata:              md,
		SkipAnalysis:          skipAnalysis,
	}
	err = artemis_orchestrations.InsertAiWorkflowAnalysisResult(ctx, wr)
	if err != nil {
		log.Err(err).Interface("filteredJsonResponses", filteredJsonResponses).Interface("jro", jro).Interface("wr", wr).Msg("UpdateTaskOutput: failed")
		return nil, err
	}
	return res, nil
}

func (z *ZeusAiPlatformActivities) RecordCompletionResponse(ctx context.Context, ou org_users.OrgUser, resp *ChatCompletionQueryResponse) (int, error) {
	if resp == nil {
		return 0, nil
	}
	b, err := json.Marshal(resp.Prompt)
	if err != nil {
		log.Err(err).Msg("RecordCompletionResponse: failed")
		return 0, err
	}
	rid, err := hera_openai_dbmodels.InsertCompletionResponseChatGpt(ctx, ou, resp.Response, b)
	if err != nil {
		log.Err(err).Msg("ZeusAiPlatformActivities: RecordCompletionResponse: failed")
		return rid, err
	}
	return rid, nil
}
