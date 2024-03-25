package ai_platform_service_orchestrations

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
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
	readOnlyFormat                      = "read-only"
	jsonFormat                          = "json"
	twitterPlatform                     = "twitter"
	redditPlatform                      = "reddit"
	discordPlatform                     = "discord"
	telegramPlatform                    = "telegram"
	webPlatform                         = "web"
)

func UnmarshallOpenAiJsonInterface(fn string, cr *ChatCompletionQueryResponse) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for _, cho := range cr.Response.Choices {
		if cho.Message.ToolCalls == nil && cho.Message.Content != "" {
			err := json.Unmarshal([]byte(cho.Message.Content), &m)
			if err != nil {
				log.Err(err).Interface("tool_calls", cho.Message.ToolCalls).Interface("cho.Message.Content", cho.Message.Content).Msg("failed to unmarshal json")
				return nil, err
			}
			for k, v := range m {
				if k == "tool_uses" {
					toolUses := v.([]interface{})
					for _, tu := range toolUses {
						tuMap := tu.(map[string]interface{})
						for k1, v1 := range tuMap {
							if k1 == "parameters" {
								tool, ok := v1.(map[string]interface{})
								if ok {
									return tool, nil
								}
								log.Info().Interface("tool", tool).Msg("tool")
							}
						}
					}
				}
			}

		}
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
	emsg, ok := m["error"]
	if ok {
		return nil, fmt.Errorf("error: %v", emsg)
	}
	return m, nil
}

func UnmarshallOpenAiJsonInterfaceSlice(fn string, cr *ChatCompletionQueryResponse) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	for _, cho := range cr.Response.Choices {
		if cho.Message.ToolCalls == nil && cho.Message.Content != "" {
			m := make(map[string]interface{})
			err := json.Unmarshal([]byte(cho.Message.Content), &m)
			if err != nil {
				log.Err(err).Interface("tool_calls", cho.Message.ToolCalls).Interface("cho.Message.Content", cho.Message.Content).Msg("failed to unmarshal json")
				return nil, err
			}
			for k, v := range m {
				if k == "tool_uses" {
					toolUses := v.([]interface{})
					for _, tu := range toolUses {
						tuMap := tu.(map[string]interface{})
						for k1, v1 := range tuMap {
							if k1 == "parameters" {
								tool, ok := v1.(map[string]interface{})
								if ok {
									results = append(results, m)
								}
								log.Info().Interface("tool", tool).Msg("tool")
							}
						}
					}
				}
			}
		}
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
