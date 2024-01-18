package ai_platform_service_orchestrations

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
)

const (
	keepTweetDescriptionField    = "add the ids from the msg_id fields you want to keep to this list"
	discardTweetDescriptionField = "add the ids from the msg_id fields you want to discard to this list"
)

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
	for _, v := range sg.SearchResults {
		msgMap[v.UnixTimestamp] = true
	}
	seen := make(map[int]bool)
	for _, v := range resp.FilteredMessages.MsgKeepIds {
		msgID, mrr := strconv.Atoi(v)
		if mrr != nil {
			log.Err(mrr).Msg("ExtractTweets: failed to convert msg id to int")
			return nil, mrr
		}
		if _, ok := seen[msgID]; ok {
			log.Info().Msgf("ExtractTweets: msgID %d already seen", msgID)
			return nil, fmt.Errorf("ExtractTweets: msgID %d already seen", msgID)
		}
		if _, ok := msgMap[msgID]; !ok {
			log.Info().Msgf("ExtractTweets: msgID %d not found in original search results", msgID)
			return nil, fmt.Errorf("ExtractTweets: msgID %d not found in original search results", msgID)
		}
		seen[msgID] = true
	}
	log.Info().Int("kept", len(resp.FilteredMessages.MsgKeepIds)).Int("all", len(msgMap)).Msg("ExtractTweets: kept messages")
	prompt := make(map[string]string)
	prompt["prompt"] = sg.ExtractionPromptExt
	resp.Prompt = prompt
	return resp, nil
}

const (
	ingestedMessageID   = "msg_id"
	ingestedMessageBody = "msg_body"
	keepMsgIDs          = "msg_keep_ids"
	discardMsgIDs       = "msg_discard_ids"
)

func FilterAndExtractRelevantTweetsJsonSchemaFunctionDef(keepMsgInst string) openai.FunctionDefinition {
	properties := make(map[string]jsonschema.Definition)
	keepMsgs := jsonschema.Definition{
		Type:        jsonschema.Array,
		Description: keepMsgInst,
		Items: &jsonschema.Definition{
			Type: jsonschema.String,
		},
	}
	properties[keepMsgIDs] = keepMsgs
	fdSchema := jsonschema.Definition{
		Type:       jsonschema.Object,
		Properties: properties,
		Required:   []string{keepMsgIDs},
	}
	fd := openai.FunctionDefinition{
		Name:       "twitter_extract_tweets",
		Parameters: fdSchema,
	}
	return fd
}
