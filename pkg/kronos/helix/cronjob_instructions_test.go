package kronos_helix

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
)

func (t *KronosWorkerTestSuite) TestCronJobWorkflowStep() {
	ojs, jerr := artemis_orchestrations.SelectSystemOrchestrationsWithInstructionsByGroup(ctx, internalOrgID, olympus)
	t.Require().Nil(jerr)
	count := 0
	for _, ojob := range ojs {
		fmt.Println(ojob.Type, ojob.GroupName)
		if ojob.Type != Cronjob {
			continue
		}
		if ojob.GroupName == olympus && ojob.Type == Cronjob {
			count++
		}
	}
	t.Require().Equal(1, count)
}

// You can change any params for this, it is a template of the other test meant for creating alerts
func (t *KronosWorkerTestSuite) TestInsertCronJobScratchPad() {
	groupName := "IrisPlatformServiceWorkflows"
	instType := "Cronjob"

	orchName := fmt.Sprintf("%s-%s", groupName, instType)
	inst := Instructions{
		GroupName: groupName,
		Type:      instType,
		CronJob: CronJobInstructions{
			Endpoint:     path.Join("https://iris.zeus.fyi/v1/internal/", "/router/serverless/refresh"),
			PollInterval: 5 * time.Minute,
		},
	}
	b, err := json.Marshal(inst)
	t.Require().Nil(err)
	groupName = olympus
	instType = Cronjob
	oj := artemis_orchestrations.OrchestrationJob{
		Orchestrations: artemis_autogen_bases.Orchestrations{
			OrgID:             t.Tc.ProductionLocalTemporalOrgID,
			Active:            true,
			GroupName:         groupName,
			Type:              instType,
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
