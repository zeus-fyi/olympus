package ai_platform_service_orchestrations

import (
	"fmt"
	"strings"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
)

// S3SaveWorkflowIO

func (t *ZeusWorkerTestSuite) TestS3HelperDownloadWfData() {
	//apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	//act := NewZeusAiPlatformActivities()
	//
	//wfr, err := act.SelectWorkflowIO(ctx, 1711583373616565000)
	//t.Require().Nil(err)
	//t.Require().NotNil(wfr)
}

func (t *ZeusWorkerTestSuite) TestSecretsSelect() {
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, org_users.NewOrgUserWithID(FlowsOrgID, 0), "s3-ovh-us-west-or")
	t.Require().Nil(err)
	t.Require().NotNil(ps)
	t.Assert().NotEmpty(ps.S3AccessKey)
	t.Assert().NotEmpty(ps.S3SecretKey)
}

func (t *ZeusWorkerTestSuite) TestS3HelperUploadWfData() {
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
	wfiNew, err := s3ws(ctx, tmp)
	t.Require().Nil(err)
	t.Require().NotNil(wfiNew)

	/*
		athena.OvhS3Manager, err = s3base.NewOvhConnS3ClientWithStaticCreds(ctx, ps.S3AccessKey, ps.S3SecretKey)
		if err != nil {
			log.Err(err).Msg("s3ws: failed to save workflow io")
			return nil, err
		}
	*/
}
