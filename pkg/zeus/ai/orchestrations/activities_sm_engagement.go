package ai_platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
)

func (z *ZeusAiPlatformActivities) AnalyzeEngagementTweets(ctx context.Context, ou org_users.OrgUser, sg *hera_search.SearchResultGroup) (*ChatCompletionQueryResponse, error) {
	if sg == nil || sg.SearchResults == nil || len(sg.SearchResults) == 0 {
		log.Info().Msg("AnalyzeEngagementTweets: no search results to analyze engagement")
		return nil, nil
	}
	za := NewZeusAiPlatformActivities()
	params := hera_openai.OpenAIParams{
		Model:              sg.Model,
		Prompt:             hera_search.FormatSearchResultsV3(sg.SearchResults),
		FunctionDefinition: sg.FunctionDefinition,
	}
	resp, err := za.CreateJsonOutputModelResponse(ctx, ou, params)
	if err != nil {
		log.Err(err).Msg("AnalyzeEngagementTweets: CreateJsonOutputModelResponse failed")
		return nil, err
	}
	var m any
	if len(resp.Response.Choices) > 0 && len(resp.Response.Choices[0].Message.ToolCalls) > 0 {
		m, err = UnmarshallOpenAiJsonInterfaceSlice(sg.FunctionDefinition.Name, resp)
		if err != nil {
			log.Err(err).Interface("m", m).Msg("UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterface failed")
			return nil, err
		}
	} else {
		m, err = UnmarshallOpenAiJsonInterface(sg.FunctionDefinition.Name, resp)
		if err != nil {
			log.Err(err).Interface("m", m).Msg("UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterface failed")
			return nil, err
		}
	}
	jsd := artemis_orchestrations.ConvertToJsonSchema(sg.FunctionDefinition)
	resp.JsonResponseResults, err = artemis_orchestrations.AssignMapValuesMultipleJsonSchemasSlice(jsd, m)
	if err != nil {
		log.Err(err).Msg("AnalyzeEngagementTweets: AssignMapValuesMultipleJsonSchemasSlice failed")
		return nil, err
	}
	return resp, nil
}
