package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/openai"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
)

func (z *ZeusAiPlatformActivities) CreateJsonOutputModelResponse(ctx context.Context, mb *MbChildSubProcessParams, params hera_openai.OpenAIParams) (*ChatCompletionQueryResponse, error) {
	tv, err := artemis_orchestrations.SelectTask(ctx, mb.Ou, mb.Tc.TaskID)
	if err != nil {
		log.Err(err).Msg("SelectTaskDefinition: failed to get task definition")
		return nil, err
	}
	if len(tv) == 0 {
		err = fmt.Errorf("failed to get task definition for task id: %d", mb.Tc.TaskID)
		log.Err(err).Msg("SelectTaskDefinition: failed to get task definition")
		return nil, err
	}
	var jsd []*artemis_orchestrations.JsonSchemaDefinition
	for _, taskDef := range tv {
		jsd = append(jsd, taskDef.Schemas...)
	}
	params.FunctionDefinition = artemis_orchestrations.ConvertToFuncDef(jsd)
	in, err := gws(ctx, mb.Wsr.InputID)
	if err != nil {
		log.Err(err).Msg("CreateJsonOutputModelResponse: gws failed")
		return nil, err
	}
	pr := in.WorkflowStageInfo.PromptReduction
	var sg *hera_search.SearchResultGroup
	if pr.PromptReductionSearchResults != nil && pr.PromptReductionSearchResults.OutSearchGroups != nil && mb.Wsr.ChunkOffset < len(pr.PromptReductionSearchResults.OutSearchGroups) {
		sg = pr.PromptReductionSearchResults.OutSearchGroups[mb.Wsr.ChunkOffset]
	} else {
		sg = &hera_search.SearchResultGroup{
			BodyPrompt:    pr.PromptReductionText.OutPromptChunks[mb.Wsr.ChunkOffset],
			SearchResults: []hera_search.SearchResult{},
		}
	}
	params.Prompt = sg.GetPromptBody()
	var resp openai.ChatCompletionResponse
	ps, err := GetMockingBirdSecrets(ctx, mb.Ou)
	if err != nil || ps == nil || ps.ApiKey == "" {
		log.Info().Msg("CreatJsonOutputModelResponse: GetMockingBirdSecrets failed to find user openai api key, using system key")
		err = nil
		resp, err = hera_openai.HeraOpenAI.MakeCodeGenRequestJsonFormattedOutput(ctx, mb.Ou, params)
	} else {
		oc := hera_openai.InitOrgHeraOpenAI(ps.ApiKey)
		resp, err = oc.MakeCodeGenRequestJsonFormattedOutput(ctx, mb.Ou, params)
	}
	if err != nil {
		log.Err(err).Interface("params", params).Msg("CreatJsonOutputModelResponse: MakeCodeGenRequestJsonFormattedOutput failed")
		return nil, err
	}

	b, err := json.Marshal(params.Prompt)
	if err != nil {
		log.Err(err).Msg("RecordCompletionResponse: failed")
		return nil, err
	}
	rid, err := hera_openai_dbmodels.InsertCompletionResponseChatGpt(ctx, mb.Ou, resp, b)
	if err != nil {
		log.Err(err).Msg("ZeusAiPlatformActivities: RecordCompletionResponse: failed")
		return nil, err
	}
	cr := &ChatCompletionQueryResponse{
		Params:     params,
		Prompt:     map[string]string{"prompt": params.Prompt},
		Response:   resp,
		ResponseID: rid,
		Schemas:    jsd,
	}
	var m any
	var anyErr error
	if len(cr.Response.Choices) > 0 && len(cr.Response.Choices[0].Message.ToolCalls) > 0 {
		m, anyErr = UnmarshallOpenAiJsonInterfaceSlice(params.FunctionDefinition.Name, cr)
		if anyErr != nil {
			log.Err(anyErr).Interface("m", m).Interface("resp", cr.Response.Choices).Msg("1_UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterfaceSlice failed")
			return nil, anyErr
		}
	} else {
		m, anyErr = UnmarshallOpenAiJsonInterface(params.FunctionDefinition.Name, cr)
		if anyErr != nil {
			log.Err(anyErr).Interface("m", m).Interface("resp", cr.Response.Choices).Msg("2_UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterface failed")
			return nil, anyErr
		}
	}
	var tmpResp []artemis_orchestrations.JsonSchemaDefinition
	tmpResp, anyErr = artemis_orchestrations.AssignMapValuesMultipleJsonSchemasSlice(jsd, m)
	if anyErr != nil {
		log.Err(anyErr).Interface("m", m).Interface("jsd", jsd).Msg("AssignMapValuesMultipleJsonSchemasSlice: UnmarshallOpenAiJsonInterface failed")
		return nil, anyErr
	}
	if tmpResp != nil {
		cr.JsonResponseResults = tmpResp
	}
	return cr, nil
}
