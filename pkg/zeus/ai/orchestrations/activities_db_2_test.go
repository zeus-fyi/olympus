package ai_platform_service_orchestrations

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (t *ZeusWorkerTestSuite) TestSelectWorkflowIO() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	act := NewZeusAiPlatformActivities()

	wfr, err := act.SelectWorkflowIO(ctx, 1711583373616565000)
	t.Require().Nil(err)
	t.Require().NotNil(wfr)
}

func (t *ZeusWorkerTestSuite) TestInsertWorkflowIO() {
	apps.Pg.InitPG(ctx, t.Tc.LocalDbPgconn)
	act := NewZeusAiPlatformActivities()

	wfr, err := act.SelectWorkflowIO(ctx, 2)
	t.Require().Nil(err)
	t.Require().NotNil(wfr)

	wfr.Logs = []string{"test", "test2"}
	tmp := &WorkflowStageIO{
		WorkflowStageReference: wfr.WorkflowStageReference,
		WorkflowStageInfo:      WorkflowStageInfo{
			//TaskToExecute: &TaskToExecute{
			//	WfID: "childWfID3",
			//	Ou:   t.Ou,
			//	Ec:   artemis_orchestrations.EvalContext{},
			//	Tc:   TaskContext{},
			//	Wft:  artemis_orchestrations.WorkflowTemplateData{},
			//},
		},
	}
	wfi, err := act.SaveWorkflowIO(ctx, tmp)
	t.Require().Nil(err)
	t.Require().NotNil(wfi)
	t.Require().NotZero(wfi.InputID)

	t.Require().Nil(err)
	wflu, err := act.SelectWorkflowIO(ctx, 2)
	t.Require().Nil(err)
	t.Require().NotNil(wflu)
	//t.Require().Equal(t.Ou, wflu.WorkflowStageInfo.TaskToExecute.Ou)
	//t.Require().Equal("childWfID3", wflu.WorkflowStageInfo.TaskToExecute.WfID)
	t.Require().NotEmpty(wflu.Logs)

	fmt.Println(strings.Join(wflu.Logs, ","))

	tmp.InputID = 0
	tmp.ChildWfID = uuid.New().String()
	wfiNew, err := act.SaveWorkflowIO(ctx, tmp)
	t.Require().Nil(err)
	t.Require().NotNil(wfiNew)
	t.Require().NotZero(wfiNew.InputID)

}
