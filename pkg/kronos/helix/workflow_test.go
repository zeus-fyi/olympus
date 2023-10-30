package kronos_helix

import (
	"encoding/json"
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

func (t *KronosWorkerTestSuite) TestWorkflowStep() {
	ojs, jerr := artemis_orchestrations.SelectSystemOrchestrationsWithInstructionsByGroup(ctx, internalOrgID, "olympus")
	t.Require().Nil(jerr)
	for _, ojob := range ojs {
		ins := Instructions{}
		err := json.Unmarshal([]byte(ojob.Instructions), &ins)
		t.Require().Nil(jerr)

		ojsFound, err := artemis_orchestrations.SelectActiveOrchestrationsWithInstructionsUsingTimeWindow(ctx, internalOrgID, ins.Type, ins.GroupName, ins.Trigger.AlertAfterTime)
		t.Require().Nil(err)
		fmt.Println(ojsFound)
	}
}

func (t *KronosWorkerTestSuite) TestCronJobWorkflowStep() {
	ka := NewKronosActivities()
	ojs, jerr := ka.GetInternalAssignments(ctx)
	t.Require().Nil(jerr)
	count := 0
	for _, oj := range ojs {

		inst, err := ka.GetInstructionsFromJob(ctx, oj)
		t.Require().Nil(err)
		t.Require().NotEmpty(inst)

		switch oj.Type {
		case alerts:
		case monitoring:
		case cronjob:
			count++

			//ta := t.Tc.DevTemporalAuth
			//ns := "kronos.ngb72"
			//hp := "kronos.ngb72.tmprl.cloud:7233"
			//ta.Namespace = ns
			//ta.HostPort = hp
			//InitKronosHelixWorker(ctx, ta)
			//cKronos := KronosServiceWorker.Worker.ConnectTemporalClient()
			//defer cKronos.Close()
			//KronosServiceWorker.Worker.RegisterWorker(cKronos)
			//err := KronosServiceWorker.Worker.Start()
			//t.Require().Nil(err)
			//
			//err = KronosServiceWorker.ExecuteKronosWorkflow2(ctx, inst)
			//t.Require().Nil(err)
			//t.Require().Equal(time.Minute*5, inst.CronJob.PollInterval)
		}
	}
	t.Require().Equal(1, count)
}
