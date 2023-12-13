package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	hera_reddit "github.com/zeus-fyi/olympus/pkg/hera/reddit"
)

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
