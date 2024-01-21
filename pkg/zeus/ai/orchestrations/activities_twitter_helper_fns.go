package ai_platform_service_orchestrations

import (
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func EvalFormatTweetForApiJsonSchema(formatType string) openai.FunctionDefinition {
	properties := make(map[string]jsonschema.Definition)
	required := []string{text}
	tweetBody := jsonschema.Definition{
		Type: jsonschema.String,
		Description: "Set this value with the text body for your suggested tweet response." +
			" It must follow these rules 1. The text content of a Tweet can contain up to 280" +
			" characters or Unicode glyphs. Spaces also count against this. 2." +
			" If you use emojis or glyphs, you should assume it will cost 3 characters from" +
			" your character budget.",
	}
	properties[text] = tweetBody
	if formatType == socialMediaEngagementResponseFormat {
		replyToTweetID := jsonschema.Definition{
			Type:        jsonschema.String,
			Description: "Set this value to the tweet id you are responding to",
		}
		properties[inReplyToTweetID] = replyToTweetID
		required = append(required, inReplyToTweetID)
	}
	fdSchema := jsonschema.Definition{
		Type:       jsonschema.Object,
		Properties: properties,
		Required:   required,
	}

	fd := openai.FunctionDefinition{
		Name:       "format_tweet_for_api",
		Parameters: fdSchema,
	}
	return fd
}
