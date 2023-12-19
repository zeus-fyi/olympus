package ai_platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (z *ZeusAiPlatformActivities) EvalLookup(ctx context.Context, ou org_users.OrgUser, evalID int) ([]artemis_orchestrations.EvalFn, error) {
	evalFn, err := artemis_orchestrations.SelectEvalFnsByOrgIDAndID(ctx, ou, evalID)
	if err != nil {
		log.Err(err).Msg("EvalLookup: failed to get eval fn")
		return nil, err
	}
	return evalFn, nil
}

func (z *ZeusAiPlatformActivities) EvalModelScoredJsonOutput(ctx context.Context, evalFn *artemis_orchestrations.EvalFn) (*artemis_orchestrations.EvalFn, error) {
	return evalFn, nil
}

func (z *ZeusAiPlatformActivities) EvalScoredViaApi(ctx context.Context, evalFn *artemis_orchestrations.EvalFn) (*artemis_orchestrations.EvalFn, error) {
	return evalFn, nil
}
