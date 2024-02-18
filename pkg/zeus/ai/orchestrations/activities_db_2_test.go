package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

func (t *ZeusWorkerTestSuite) TestSelectWorkflowIO() {
	apps.Pg.InitPG(ctx, t.Tc.LocalDbPgconn)
	act := NewZeusAiPlatformActivities()

	_, err := act.SelectWorkflowIO(ctx, 1)
	t.Require().Nil(err)
}

func (t *ZeusWorkerTestSuite) TestInsertWorkflowIO() {
	apps.Pg.InitPG(ctx, t.Tc.LocalDbPgconn)
	act := NewZeusAiPlatformActivities()

	wfr, err := act.SelectWorkflowIO(ctx, 1)
	t.Require().Nil(err)
	t.Require().NotNil(wfr)
	err = act.SaveWorkflowIO(ctx, &WorkflowStageIO{
		WorkflowStageReference: wfr.WorkflowStageReference,
		WorkflowStageInfo: WorkflowStageInfo{
			TaskToExecute: &TaskToExecute{
				WfID:        "",
				Ou:          t.Ou,
				Ec:          artemis_orchestrations.EvalContext{},
				Tc:          TaskContext{},
				Wft:         artemis_orchestrations.WorkflowTemplateData{},
				Sg:          nil,
				Wr:          nil,
				RetryPolicy: nil,
			},
		},
	})
	t.Require().Nil(err)
	wflu, err := act.SelectWorkflowIO(ctx, 1)
	t.Require().Nil(err)
	t.Require().NotNil(wflu)
	t.Require().Equal(t.Ou, wflu.WorkflowStageInfo.TaskToExecute.Ou)
}
