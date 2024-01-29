package ai_platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (z *ZeusAiPlatformActivities) SelectRetrievalTask(ctx context.Context, ou org_users.OrgUser, retID int) ([]artemis_orchestrations.RetrievalItem, error) {
	resp, err := artemis_orchestrations.SelectRetrievals(ctx, ou, retID)
	if err != nil {
		log.Err(err).Interface("resp", resp).Msg("SaveTaskOutput: failed")
		return resp, err
	}
	return resp, nil
}
