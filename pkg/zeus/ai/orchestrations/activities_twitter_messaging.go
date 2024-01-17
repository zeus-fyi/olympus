package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	hera_twitter "github.com/zeus-fyi/olympus/pkg/hera/twitter"
)

const (
	keepTweetDescriptionField2    = "add the ids from the msg_id fields you want to keep to this list"
	discardTweetDescriptionField2 = "add the ids from the msg_id fields you want to discard to this list"
)

func (z *ZeusAiPlatformActivities) CreateTweets(ctx context.Context, ou org_users.OrgUser, sg *hera_search.SearchResultGroup) (*ChatCompletionQueryResponse, error) {
	if sg == nil || sg.ExtractionPromptExt == "" {
		log.Info().Msg("CreateTweets: no prompt was provided")
		return nil, nil
	}

	isReply := false
	if sg.ResponseFormat != socialMediaEngagementResponseFormat {
		isReply = true
	}
	fd := TwitterTweetCreationApiCallJsonSchemaFunctionDef(isReply)
	za := NewZeusAiPlatformActivities()
	resp, err := za.CreateJsonOutputModelResponse(ctx, ou, hera_openai.OpenAIParams{
		Model:              sg.Model,
		FunctionDefinition: fd,
		Prompt:             sg.ExtractionPromptExt,
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

type Tweets struct {
	Outbound []TweetResponse `json:"keep_tweet_ids"`
}

type TweetResponse struct {
	Text         string `json:"text"`
	ReplyTweetID string `json:"in_reply_to_tweet_id"`
}

func (z *ZeusAiPlatformActivities) SocialTweetTask(ctx context.Context, ou org_users.OrgUser, tr *TweetResponse) (*TweetResponse, error) {
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, twitterPlatform)
	if err != nil || ps == nil || ps.ApiKey == "" {
		if err == nil {
			err = fmt.Errorf("failed to get mockingbird secrets")
		}
		log.Err(err).Msg("AiAnalysisTask: GetMockingbirdPlatformSecrets: failed to get mockingbird secrets")
		return nil, err
	}
	tc, err := hera_twitter.InitOrgTwitterClient(ctx, ps.OAuth2Public, ps.OAuth2Secret)
	if err != nil {
		log.Err(err).Msg("TweetTask: InitTwitterClient failed")
		return nil, err
	}
	tweet, err := tc.V2alt.CreateTweet(ctx, twitter.CreateTweetRequest{
		Text:  tr.Text,
		Reply: &twitter.CreateTweetReply{InReplyToTweetID: tr.ReplyTweetID},
	})
	if err != nil {
		log.Err(err).Msg("TweetTask: CreateTweet failed")
		return nil, err
	}
	log.Info().Msgf("TweetTask: tweet created: %v", tweet)
	return nil, nil
}

const (
	replyToDescription = "a list of message ids to reply to"
	replyToIds         = "in_reply_to_tweet_id"
)

func TwitterTweetCreationApiCallJsonSchemaFunctionDef(isReplyTo bool) openai.FunctionDefinition {
	required := []string{text}
	if isReplyTo {
		required = append(required, replyToIds)
	}
	tweetDefinition := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			text: {
				Type: jsonschema.Array,
				Items: &jsonschema.Definition{
					Type: jsonschema.String,
				},
			},
			replyToIds: {
				Type: jsonschema.Array,
				Items: &jsonschema.Definition{
					Type: jsonschema.String,
				},
			},
		},
		Required: required,
	}
	properties := make(map[string]jsonschema.Definition)
	keepMsgs := jsonschema.Definition{
		Type:        jsonschema.Array,
		Description: replyToDescription,
		Items:       &tweetDefinition,
		Properties:  properties,
	}
	fd := openai.FunctionDefinition{
		Name:       "twitter_create_tweet_suggestions",
		Parameters: keepMsgs,
	}
	return fd
}
