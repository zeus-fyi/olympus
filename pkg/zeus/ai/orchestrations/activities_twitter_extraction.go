package ai_platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
)

const (
	keepTweetDescriptionField = "add the msg_id with the msg_body field that you want to keep to this array of msg_ids"
)

type SocialExtractionParams struct {
	Sg   *hera_search.SearchResultGroup
	Resp *ChatCompletionQueryResponse
}

func (z *ZeusAiPlatformActivities) ExtractTweets(ctx context.Context, ou org_users.OrgUser, sg *hera_search.SearchResultGroup) (*ChatCompletionQueryResponse, error) {
	if sg == nil || sg.SearchResults == nil || len(sg.SearchResults) == 0 {
		log.Info().Msg("ExtractTweets: no search results")
		return nil, nil
	}
	fd := FilterAndExtractRelevantTweetsJsonSchemaFunctionDef(keepTweetDescriptionField)
	za := NewZeusAiPlatformActivities()
	params := hera_openai.OpenAIParams{
		Model:              sg.Model,
		FunctionDefinition: fd,
		Prompt:             hera_search.FormatSearchResultsV3(sg.SearchResults),
		SystemPromptExt:    sg.ExtractionPromptExt,
	}
	resp, err := za.CreateJsonOutputModelResponse(ctx, ou, params)
	if err != nil {
		log.Err(err).Msg("ExtractTweets: CreateJsonOutputModelResponse failed")
		return nil, err
	}
	err = UnmarshallFilteredMsgIdsFromAiJson(fd.Name, resp)
	if err != nil {
		log.Err(err).Msg("ExtractTweets: failed to unmarshall json interface")
		return nil, err
	}

	msgMap := make(map[int]bool)
	srMap := make(map[int]hera_search.SearchResult)
	for _, v := range sg.SearchResults {
		msgMap[v.UnixTimestamp] = true
		srMap[v.UnixTimestamp] = v
	}
	seen := make(map[int]bool)
	for _, msgID := range resp.FilteredMessages.MsgKeepIds {
		if _, ok := seen[msgID]; ok {
			log.Warn().Msgf("ExtractTweets: msgID %d already seen", msgID)
			continue
		}
		if _, ok := msgMap[msgID]; !ok {
			log.Info().Msgf("ExtractTweets: msgID %d not found in original search results", msgID)
			continue
		}
		resp.FilteredSearchResults = append(resp.FilteredSearchResults, srMap[msgID])
		seen[msgID] = true
	}
	log.Info().Int("kept", len(resp.FilteredMessages.MsgKeepIds)).Int("all", len(msgMap)).Msg("ExtractTweets: kept messages")
	if resp.Prompt == nil {
		resp.Prompt = make(map[string]string)
	}
	resp.Prompt["extraction-prompt"] = sg.ExtractionPromptExt
	return resp, nil
}

const (
	filteredMessageArrayJsonFieldName = "msg_ids"
	ingestedMessageBody               = "msg_body"
)

func FilterAndExtractRelevantTweetsJsonSchemaFunctionDef(keepMsgInst string) openai.FunctionDefinition {
	properties := make(map[string]jsonschema.Definition)
	keepMsgs := jsonschema.Definition{
		Type:        jsonschema.Array,
		Description: keepMsgInst,
		Items: &jsonschema.Definition{
			Type: jsonschema.Integer,
		},
	}
	properties[filteredMessageArrayJsonFieldName] = keepMsgs
	fdSchema := jsonschema.Definition{
		Type:       jsonschema.Object,
		Properties: properties,
		Required:   []string{filteredMessageArrayJsonFieldName},
	}
	fd := openai.FunctionDefinition{
		Name:       "twitter_extract_tweets",
		Parameters: fdSchema,
	}
	return fd
}
