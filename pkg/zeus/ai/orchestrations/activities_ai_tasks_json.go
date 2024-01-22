package ai_platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
)

func (z *ZeusAiPlatformActivities) CreateJsonOutputModelResponse(ctx context.Context, ou org_users.OrgUser, params hera_openai.OpenAIParams) (*ChatCompletionQueryResponse, error) {
	var err error
	var resp openai.ChatCompletionResponse
	ps, err := GetMockingBirdSecrets(ctx, ou)
	if err != nil || ps == nil || ps.ApiKey == "" {
		log.Info().Msg("CreatJsonOutputModelResponse: GetMockingBirdSecrets failed to find user openai api key, using system key")
		err = nil
		resp, err = hera_openai.HeraOpenAI.MakeCodeGenRequestJsonFormattedOutput(ctx, ou, params)
	} else {
		oc := hera_openai.InitOrgHeraOpenAI(ps.ApiKey)
		resp, err = oc.MakeCodeGenRequestJsonFormattedOutput(ctx, ou, params)
	}
	if err != nil {
		log.Err(err).Msg("CreatJsonOutputModelResponse: MakeCodeGenRequestJsonFormattedOutput failed")
		return nil, err
	}
	cr := &ChatCompletionQueryResponse{
		Prompt:   map[string]string{"prompt": params.Prompt},
		Response: resp,
	}
	var m any
	if len(cr.Response.Choices) > 0 && len(cr.Response.Choices[0].Message.ToolCalls) > 0 {
		m, err = UnmarshallOpenAiJsonInterfaceSlice(params.FunctionDefinition.Name, cr)
		if err != nil {
			log.Err(err).Interface("m", m).Msg("UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterfaceSlice failed")
			return nil, err
		}

	} else {
		m, err = UnmarshallOpenAiJsonInterface(params.FunctionDefinition.Name, cr)
		if err != nil {
			log.Err(err).Interface("m", m).Msg("UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterface failed")
			return nil, err
		}
	}
	jsd := artemis_orchestrations.ConvertToJsonSchema(params.FunctionDefinition)
	cr.JsonResponseResults = artemis_orchestrations.AssignMapValuesMultipleJsonSchemasSlice(jsd, m)
	return cr, nil
}
