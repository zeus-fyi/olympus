package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
)

const (
	Gpt4JsonModel = "gpt-4-1106-preview"
	Gpt3JsonModel = "gpt-3.5-turbo-1106"
)

const (
	text                                = "text"
	inReplyToTweetID                    = "in_reply_to_tweet_id"
	socialMediaEngagementResponseFormat = "social-media-engagement"
	socialMediaExtractionResponseFormat = "social-media-extraction"
	jsonFormat                          = "json"
	twitterPlatform                     = "twitter"
	redditPlatform                      = "reddit"
	discordPlatform                     = "discord"
	telegramPlatform                    = "telegram"
	webPlatform                         = "web"
)

func (z *ZeusAiPlatformActivities) EvalFormatForApi(ctx context.Context, ou org_users.OrgUser, ta *artemis_orchestrations.TriggerAction) (*ChatCompletionQueryResponse, error) {
	var fnApiFormat openai.FunctionDefinition

	platformName := ""
	switch platformName {
	case twitterPlatform:
		fnApiFormat = EvalFormatTweetForApiJsonSchema(ta.TriggerAction)
		log.Info().Msgf("EvalFormatTweetForApi: body: %v", fnApiFormat)
	default:
		return nil, fmt.Errorf("EvalFormatForApi: platform %s not supported", platformName)
	}
	za := NewZeusAiPlatformActivities()
	resp, err := za.CreateJsonOutputModelResponse(ctx, ou, hera_openai.OpenAIParams{
		Model:              Gpt4JsonModel,
		FunctionDefinition: fnApiFormat,
	})
	if err != nil {
		log.Err(err).Msg("EvalFormatTweetForApi: CreateJsonOutputModelResponse failed")
		return nil, err
	}
	switch platformName {
	case twitterPlatform:
		tr, terr := UnmarshallTwitterFromAiJson(fnApiFormat.Name, resp)
		if terr != nil {
			log.Err(terr).Msg("EvalFormatTweetForApi: failed to unmarshall json interface")
			return nil, terr
		}
		resp.CreateTweetRequest = tr
	case redditPlatform:
	case discordPlatform:
	case telegramPlatform:
	default:
		return nil, fmt.Errorf("EvalFormatForApi: platform %s not supported", platformName)
	}
	return resp, nil
}

func UnmarshallTwitterFromAiJson(fn string, cr *ChatCompletionQueryResponse) (*twitter.CreateTweetRequest, error) {
	m, err := UnmarshallOpenAiJsonInterface(fn, cr)
	if err != nil {
		return nil, err
	}
	rb, ok := m[text].(string)
	if !ok || len(rb) == 0 {
		return nil, fmt.Errorf("text body had no text, or was not a string")
	}
	tr := &twitter.CreateTweetRequest{
		Text: rb,
	}
	rt, ok := m[inReplyToTweetID].(string)
	if ok {
		tr.Reply = &twitter.CreateTweetReply{InReplyToTweetID: rt}
	}
	return tr, nil
}

func UnmarshallOpenAiJsonInterface(fn string, cr *ChatCompletionQueryResponse) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for _, cho := range cr.Response.Choices {
		for _, tvr := range cho.Message.ToolCalls {
			if tvr.Function.Name == fn {
				err := json.Unmarshal([]byte(tvr.Function.Arguments), &m)
				if err != nil {
					log.Err(err).Interface("tool_calls", cho.Message.ToolCalls).Interface("tvr", tvr).Msg("failed to unmarshal json")
					return nil, err
				}
			}
		}
	}
	return m, nil
}

func UnmarshallOpenAiJsonInterfaceSlice(fn string, cr *ChatCompletionQueryResponse) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	for _, cho := range cr.Response.Choices {
		for _, tvr := range cho.Message.ToolCalls {
			if tvr.Function.Name == fn {
				m := make(map[string]interface{})
				err := json.Unmarshal([]byte(tvr.Function.Arguments), &m)
				if err != nil {
					log.Err(err).Interface("results", results).Interface("tvr", tvr).Msg("failed to unmarshal json")
					return nil, err
				}
				results = append(results, m)
			}
		}
	}
	return results, nil
}
