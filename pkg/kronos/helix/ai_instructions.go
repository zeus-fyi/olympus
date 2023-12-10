package kronos_helix

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func InsertSearchIndexerOrchestration(ctx context.Context, ou org_users.OrgUser, platform string) error {
	inst := Instructions{
		GroupName: mockingbird,
		Type:      cronjob,
		CronJob: CronJobInstructions{
			PollInterval: 5 * time.Minute,
		},
	}
	orchName := fmt.Sprintf("Mockingbird-Search-Indexer-Cronjob-%s", platform)
	b, err := json.Marshal(inst)
	if err != nil {
		return err
	}
	oj := artemis_orchestrations.OrchestrationJob{
		Orchestrations: artemis_autogen_bases.Orchestrations{
			OrgID:             ou.OrgID,
			Active:            true,
			GroupName:         mockingbird,
			Type:              cronjob,
			Instructions:      string(b),
			OrchestrationName: orchName,
		},
	}
	err = oj.UpsertOrchestrationWithInstructions(ctx)
	if err != nil {
		log.Err(err).Interface("ou", ou).Interface("orchName", orchName).Msg("InsertSearchIndexerOrchestration: UpsertOrchestrationWithInstructions failed")
		return err
	}
	return nil
}
