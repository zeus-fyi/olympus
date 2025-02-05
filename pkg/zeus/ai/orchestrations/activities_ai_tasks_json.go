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
	"go.temporal.io/sdk/activity"
)

const (
	FlowsOrgID      = 1685378241971196000
	FlowsS3Ovh      = "s3-ovh-us-west-or"
	FlowsBucketName = "flows"
)

func (z *ZeusAiPlatformActivities) CreateJsonOutputModelResponse(ctx context.Context, mb *MbChildSubProcessParams, params hera_openai.OpenAIParams) (*ChatCompletionQueryResponse, error) {
	jsd, err := getJsonSchemaDefs(ctx, mb)
	if err != nil {
		log.Err(err).Msg("CreateJsonOutputModelResponse: getJsonSchemaDefs failed")
		return nil, err
	}
	// todo: needs to save data outputs
	params.FunctionDefinition = artemis_orchestrations.ConvertToFuncDef(jsd)
	in, err := gs3wfs(ctx, mb)
	if err != nil {
		log.Err(err).Msg("CreateJsonOutputModelResponse: gws failed")
		return nil, err
	}
	sg := getJsonSgChunkToProcess(mb, in)
	if mb.WfExecParams.WorkflowOverrides.TaskPromptOverrides != nil {
		if v, ok := mb.WfExecParams.WorkflowOverrides.TaskPromptOverrides[mb.Tc.TaskName]; ok {
			if len(v.SystemPromptExt) > 0 {
				params.SystemPromptExt = v.SystemPromptExt
				log.Info().Interface("v.SystemPromptExt", v).Msg("CreateJsonOutputModelResponse")
			}
		}
	}

	params.Prompt = sg.GetPromptBody()
	if len(params.Prompt) == 0 {
		log.Warn().Interface("mb.Tc.TaskName", mb.Tc.TaskName).Msg("CreateJsonOutputModelResponse: prompt is empty")
		return nil, fmt.Errorf("CreateJsonOutputModelResponse: prompt is empty")
	}
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
	jsv, err := unmarshallJsonOpenAI(cr, jsd, sg)
	if err != nil {
		log.Err(err).Interface("jsv", jsv).Msg("unmarshallJsonOpenAI failed")
	}
	if len(jsv) > 0 {
		cr.JsonResponseResults = jsv
	}
	log.Info().Interface("len(cr.JsonResponseResults)", len(cr.JsonResponseResults)).Interface("len(sg.RegexSearchResults)", len(sg.RegexSearchResults)).Interface("len(sg.ApiResponseResults)", len(sg.ApiResponseResults)).Msg("CreateJsonOutputModelResponse }")

	// temp
	var jsff []artemis_orchestrations.JsonSchemaDefinition
	for _, jt := range jsd {
		if jt != nil {
			jsff = append(jsff, *jt)
		}
	}
	//m := make(map[string]bool)
	//for _, v := range sg.ApiResponseResults {
	//	m[v.Source] = true
	//}
	payloadMaps := artemis_orchestrations.CreateMapInterfaceFromAssignedSchemaFields(jsff)
	for _, pl := range payloadMaps {
		fmt.Println(pl)
		//tv, ok := pl["entity"]
		//if !ok {
		//	return nil, fmt.Errorf("not ok")
		//}
		//sv, ok := tv.(string)
		//if !ok {
		//	return nil, fmt.Errorf("not ok")
		//}
		//mv, ok := m[sv]
		//if !mv || !ok {
		//	return nil, fmt.Errorf("not ok")
		//}
	}
	//dj := DebugJsonOutputs{
	//	Mb:            mb,
	//	Params:        params,
	//	JsonResponses: payloadMaps,
	//}
	//dj.Save()
	// end temp
	activity.RecordHeartbeat(ctx, cr.Response.ID)
	return cr, nil
}

func getJsonSgChunkToProcess(mb *MbChildSubProcessParams, in *WorkflowStageIO) *hera_search.SearchResultGroup {
	pr := in.WorkflowStageInfo.PromptReduction
	var sg *hera_search.SearchResultGroup
	if in.WorkflowStageInfo.PromptTextFromTextStage != "" {
		sg = &hera_search.SearchResultGroup{
			BodyPrompt: in.WorkflowStageInfo.PromptTextFromTextStage,
		}
	} else if pr.PromptReductionSearchResults != nil && pr.PromptReductionSearchResults.OutSearchGroups != nil && mb.Wsr.ChunkOffset < len(pr.PromptReductionSearchResults.OutSearchGroups) {
		sg = pr.PromptReductionSearchResults.OutSearchGroups[mb.Wsr.ChunkOffset]
	} else {
		sg = &hera_search.SearchResultGroup{
			BodyPrompt:    pr.PromptReductionText.OutPromptChunks[mb.Wsr.ChunkOffset],
			SearchResults: []hera_search.SearchResult{},
		}
	}
	return sg
}

func getJsonSgChunkToProcess2(chunk int, mb *MbChildSubProcessParams, in *WorkflowStageIO) *hera_search.SearchResultGroup {
	pr := in.WorkflowStageInfo.PromptReduction
	var sg *hera_search.SearchResultGroup
	if in.WorkflowStageInfo.PromptTextFromTextStage != "" {
		sg = &hera_search.SearchResultGroup{
			BodyPrompt: in.WorkflowStageInfo.PromptTextFromTextStage,
		}
	} else if pr.PromptReductionSearchResults != nil && pr.PromptReductionSearchResults.OutSearchGroups != nil && mb.Wsr.ChunkOffset < len(pr.PromptReductionSearchResults.OutSearchGroups) {
		sg = pr.PromptReductionSearchResults.OutSearchGroups[chunk]
	} else if pr.PromptReductionText.OutPromptChunks != nil {
		sg = &hera_search.SearchResultGroup{
			BodyPrompt:    pr.PromptReductionText.OutPromptChunks[chunk],
			SearchResults: []hera_search.SearchResult{},
		}
	} else if pr.PromptReductionSearchResults.InSearchGroup != nil {
		sg = &hera_search.SearchResultGroup{
			SearchResults: pr.PromptReductionSearchResults.InSearchGroup.ApiResponseResults,
		}
	}
	return sg
}

func unmarshallJsonOpenAI(cr *ChatCompletionQueryResponse, jsd []*artemis_orchestrations.JsonSchemaDefinition, sg *hera_search.SearchResultGroup) ([]artemis_orchestrations.JsonSchemaDefinition, error) {
	var m any
	var anyErr error
	if len(cr.Response.Choices) > 0 && len(cr.Response.Choices[0].Message.ToolCalls) > 0 {
		m, anyErr = UnmarshallOpenAiJsonInterfaceSlice(cr.Params.FunctionDefinition.Name, cr)
		if anyErr != nil {
			log.Err(anyErr).Interface("m", m).Interface("resp", cr.Response.Choices).Interface("prompt", sg.GetPromptBody()).Msg("1_UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterfaceSlice failed")
			return nil, anyErr
		}
	} else {
		m, anyErr = UnmarshallOpenAiJsonInterface(cr.Params.FunctionDefinition.Name, cr)
		if anyErr != nil {
			log.Err(anyErr).Interface("m", m).Interface("resp", cr.Response.Choices).Interface("prompt", sg.GetPromptBody()).Msg("2_UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterface failed")
			return nil, anyErr
		}
	}
	var tmpResp []artemis_orchestrations.JsonSchemaDefinition
	tmpResp, anyErr = artemis_orchestrations.AssignMapValuesMultipleJsonSchemasSlice(jsd, m)
	if anyErr != nil {
		log.Err(anyErr).Interface("m", m).Interface("jsd", jsd).Interface("prompt", sg.GetPromptBody()).Msg("AssignMapValuesMultipleJsonSchemasSlice: UnmarshallOpenAiJsonInterface failed")
		return nil, anyErr
	}
	return tmpResp, anyErr
}

func getJsonSchemaDefs(ctx context.Context, mb *MbChildSubProcessParams) ([]*artemis_orchestrations.JsonSchemaDefinition, error) {
	var jsd []*artemis_orchestrations.JsonSchemaDefinition
	if mb.Tc.EvalID > 0 && mb.Tc.EvalSchemas != nil && len(mb.Tc.EvalSchemas) > 0 {
		jsd = append(jsd, mb.Tc.EvalSchemas...)
	} else {
		tmpOu := mb.Ou
		if mb.WfExecParams.WorkflowOverrides.IsUsingFlows {
			tmpOu.OrgID = FlowsOrgID
		}
		tv, err := artemis_orchestrations.SelectTask(ctx, tmpOu, mb.Tc.TaskID)
		if err != nil {
			log.Err(err).Msg("SelectTaskDefinition: failed to get task definition")
			return nil, err
		}
		if len(tv) == 0 {
			err = fmt.Errorf("failed to get task definition for task id: %d", mb.Tc.TaskID)
			log.Err(err).Msg("SelectTaskDefinition: failed to get task definition")
			return nil, err
		}
		for _, taskDef := range tv {
			for _, sv := range taskDef.Schemas {
				if mb.WfExecParams.WorkflowOverrides.SchemaFieldOverrides != nil {
					if _, ok := mb.WfExecParams.WorkflowOverrides.SchemaFieldOverrides[sv.SchemaName]; !ok {
						continue
					}
				}
				var extFields []artemis_orchestrations.JsonSchemaField
				for si, sf := range sv.Fields {
					if mb.WfExecParams.WorkflowOverrides.SchemaFieldOverrides[sv.SchemaName] != nil {
						if fo, ok := mb.WfExecParams.WorkflowOverrides.SchemaFieldOverrides[sv.SchemaName][sf.FieldName]; ok {
							fn := sf.FieldName
							for oi, orv := range fo {
								if oi <= 0 {
									sv.Fields[si].FieldDescription = orv
								} else {
									sf.FieldName = fmt.Sprintf("%s_%d", fn, oi)
									sf.FieldDescription = orv
									extFields = append(extFields, sf)
								}
							}
						}
					}
				}
				sv.Fields = append(sv.Fields, extFields...)
			}
			jsd = append(jsd, taskDef.Schemas...)
		}
	}
	return jsd, nil
}
