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
	keepTweetDescriptionField    = "add the ids from the msg_id fields you want to keep to this list"
	discardTweetDescriptionField = "add the ids from the msg_id fields you want to discard to this list"
)

func (z *ZeusAiPlatformActivities) ExtractTweets(ctx context.Context, ou org_users.OrgUser, sg *hera_search.SearchResultGroup) (*ChatCompletionQueryResponse, error) {
	if sg == nil || sg.SearchResults == nil || len(sg.SearchResults) == 0 {
		log.Info().Msg("ExtractTweets: no search results")
		return nil, nil
	}
	fd := FilterAndExtractRelevantTweetsJsonSchemaFunctionDef(keepTweetDescriptionField, discardTweetDescriptionField)
	za := NewZeusAiPlatformActivities()
	resp, err := za.CreateJsonOutputModelResponse(ctx, ou, hera_openai.OpenAIParams{
		Model:              sg.Model,
		FunctionDefinition: fd,
		Prompt:             hera_search.FormatSearchResultsV3(sg.SearchResults),
		SystemPromptExt:    sg.ExtractionPromptExt,
	})
	if err != nil {
		log.Err(err).Msg("ExtractTweets: CreateJsonOutputModelResponse failed")
		return nil, err
	}
	prompt := make(map[string]string)
	prompt["prompt"] = sg.ExtractionPromptExt
	return &ChatCompletionQueryResponse{
		Prompt:   prompt,
		Response: resp.Response,
	}, nil
}

const (
	ingestedMessageID   = "msg_id"
	ingestedMessageBody = "msg_body"
	keepMsgIDs          = "msg_keep_ids"
	discardMsgIDs       = "msg_discard_ids"
)

func FilterAndExtractRelevantTweetsJsonSchemaFunctionDef(keepMsgInst, discardMsgInst string) openai.FunctionDefinition {
	properties := make(map[string]jsonschema.Definition)
	keepMsgs := jsonschema.Definition{
		Type:        jsonschema.Array,
		Description: keepMsgInst,
		Items: &jsonschema.Definition{
			Type: jsonschema.String,
		},
	}
	properties[keepMsgIDs] = keepMsgs
	//discardMsgs := jsonschema.Definition{
	//	Type:        jsonschema.Array,
	//	Description: discardMsgInst,
	//	Items: &jsonschema.Definition{
	//		Type: jsonschema.String,
	//	},
	//}
	//properties[discardMsgIDs] = discardMsgs
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
