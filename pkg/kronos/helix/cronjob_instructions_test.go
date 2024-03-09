package kronos_helix

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
)

// You can change any params for this, it is a template of the other test meant for creating alerts
func (t *KronosWorkerTestSuite) TestInsertCronJobScratchPad() {
	inst := Instructions{
		GroupName: olympus,
		Type:      cronjob,
		CronJob: CronJobInstructions{
			Endpoint:     fmt.Sprintf("https://api.zeus.fyi/v1/webhooks/twillio"),
			PollInterval: 5 * time.Minute,
		},
	}
	b, err := json.Marshal(inst)
	t.Require().Nil(err)

	groupName := "ZeusAiPlatformServiceWorkflows"
	instType := "Mockingbird-Twillio-Indexer-Cronjob-0"

	orchName := fmt.Sprintf("%s-%s", groupName, instType)
	oj := artemis_orchestrations.OrchestrationJob{
		Orchestrations: artemis_autogen_bases.Orchestrations{
			OrgID:             t.Tc.ProductionLocalTemporalOrgID,
			Active:            true,
			GroupName:         olympus,
			Type:              cronjob,
			Instructions:      string(b),
			OrchestrationName: orchName,
		},
	}
	err = oj.UpsertOrchestrationWithInstructions(ctx)
	t.Require().Nil(err)
	t.Assert().NotZero(oj.OrchestrationID)
}

func (t *KronosWorkerTestSuite) TestStartCronJobWorkflow() {
	groupName := "IrisPlatformServiceWorkflows"
	instType := "Cronjob"

	//endpoint :=  path.Join("https://iris.zeus.fyi/v1/internal/", "/router/serverless/refresh")
	artemis_orchestration_auth.Bearer = t.Tc.ProductionLocalTemporalBearerToken
	endpoint := fmt.Sprintf("http://localhost:8080/v1/internal/%s", "router/serverless/refresh")
	inst := Instructions{
		GroupName: groupName,
		Type:      instType,
		CronJob: CronJobInstructions{
			Endpoint:     endpoint,
			PollInterval: 5 * time.Minute,
		},
	}
	a := NewKronosActivities()
	err := a.StartCronJobWorkflow(ctx, inst)
	t.Require().Nil(err)
}
