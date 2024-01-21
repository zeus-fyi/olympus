package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
)

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
