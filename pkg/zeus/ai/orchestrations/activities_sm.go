package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/rs/zerolog/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	hera_reddit "github.com/zeus-fyi/olympus/pkg/hera/reddit"
)

type SocialMediaPlatformResponses struct {
	ModelResponse *ChatCompletionQueryResponse
	Twitter       *twitter.CreateTweetRequest `json:"twitter,omitempty"`
	Reddit        hera_reddit.SocialReddit    `json:"reddit,omitempty"`
}

func (z *ZeusAiPlatformActivities) SocialRedditTask(ctx context.Context, ou org_users.OrgUser, reply *SocialMediaPlatformResponses, sr []hera_search.SearchResult) (*ChatCompletionQueryResponse, error) {
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, "reddit")
	if err != nil {
		log.Err(err).Msg("SearchRedditNewPostsUsingSubreddit: failed to get mockingbird secrets")
		return nil, err
	}
	if ps == nil {
		return nil, fmt.Errorf("SearchRedditNewPostsUsingSubreddit: ps is nil")
	}
	if ps.OAuth2Public == "" || ps.OAuth2Secret == "" || ps.Username == "" || ps.Password == "" {
		return nil, fmt.Errorf("SearchRedditNewPostsUsingSubreddit: ps is empty")
	}
	rc, err := hera_reddit.InitOrgRedditClient(ctx, ps.OAuth2Public, ps.OAuth2Secret, ps.Username, ps.Password)
	if err != nil {
		log.Err(err).Msg("SearchRedditNewPostsUsingSubreddit: failed to init reddit client")
		return nil, err
	}
	cm, rr, err := rc.ReplyToPost(ctx, reply.Reddit.RedditReplyToPost.ParentID, reply.Reddit.RedditReplyToPost.Reply)
	if err != nil {
		log.Err(err).Msg("SearchRedditNewPostsUsingSubreddit: failed to init reddit client")
		return nil, err
	}
	log.Info().Msgf("SearchRedditNewPostsUsingSubreddit: cm: %v, rr: %v", cm, rr)
	return nil, err
}

func (z *ZeusAiPlatformActivities) SearchRedditNewPostsUsingSubreddit(ctx context.Context, ou org_users.OrgUser, subreddit string, lpo *reddit.ListOptions) ([]*reddit.Post, error) {
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, "reddit")
	if err != nil {
		log.Err(err).Msg("SearchRedditNewPostsUsingSubreddit: failed to get mockingbird secrets")
		return nil, err
	}
	if ps == nil {
		return nil, fmt.Errorf("SearchRedditNewPostsUsingSubreddit: ps is nil")
	}
	if ps.OAuth2Public == "" || ps.OAuth2Secret == "" || ps.Username == "" || ps.Password == "" {
		return nil, fmt.Errorf("SearchRedditNewPostsUsingSubreddit: ps is empty")
	}
	rc, err := hera_reddit.InitOrgRedditClient(ctx, ps.OAuth2Public, ps.OAuth2Secret, ps.Username, ps.Password)
	if err != nil {
		log.Err(err).Msg("SearchRedditNewPostsUsingSubreddit: failed to init reddit client")
		return nil, err
	}
	posts, err := rc.GetNewPosts(ctx, subreddit, lpo)
	if err != nil {
		log.Err(err).Interface("resp", posts).Msg("SearchRedditNewPostsUsingSubreddit")
		return nil, err
	}
	return posts, nil
}

func (z *ZeusAiPlatformActivities) SocialDiscordTask(ctx context.Context, ou org_users.OrgUser, reply *ChatCompletionQueryResponse, sr []hera_search.SearchResult) (*ChatCompletionQueryResponse, error) {
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, "discord")
	if err != nil || ps == nil || ps.ApiKey == "" {
		if err == nil {
			err = fmt.Errorf("failed to get mockingbird secrets")
		}
		log.Err(err).Msg("SocialDiscordTask: AiAnalysisTask: GetMockingbirdPlatformSecrets: failed to get mockingbird secrets")
		return nil, err
	}
	return nil, err
}

func (z *ZeusAiPlatformActivities) SocialTelegramTask(ctx context.Context, ou org_users.OrgUser, reply *SocialMediaPlatformResponses, sr []hera_search.SearchResult) (*ChatCompletionQueryResponse, error) {
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, "telegram")
	if err != nil {
		log.Err(err).Msg("SocialTelegramTask: failed to get mockingbird secrets")
		return nil, err
	}
	if ps == nil {
		return nil, fmt.Errorf("SocialTelegramTask: ps is nil")
	}
	if ps.OAuth2Public == "" || ps.OAuth2Secret == "" || ps.Username == "" || ps.Password == "" {
		return nil, fmt.Errorf("SocialTelegramTask: ps is empty")
	}
	return nil, err
}
