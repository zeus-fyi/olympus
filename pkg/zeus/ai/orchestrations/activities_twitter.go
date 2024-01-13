package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/rs/zerolog/log"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	hera_twitter "github.com/zeus-fyi/olympus/pkg/hera/twitter"
)

func (z *ZeusAiPlatformActivities) SocialTweetTask(ctx context.Context, ou org_users.OrgUser, reply *ChatCompletionQueryResponse, sr []hera_search.SearchResult) (*ChatCompletionQueryResponse, error) {
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, "twitter")
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
		Text:  "",
		Reply: &twitter.CreateTweetReply{InReplyToTweetID: ""},
	})
	if err != nil {
		log.Err(err).Msg("TweetTask: CreateTweet failed")
		return nil, err
	}
	log.Info().Msgf("TweetTask: tweet created: %v", tweet)
	return nil, nil
}
