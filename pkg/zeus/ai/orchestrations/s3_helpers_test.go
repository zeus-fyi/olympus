package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
)

// S3SaveWorkflowIO

func (t *ZeusWorkerTestSuite) TestS3HelperDownloadWfData() {
	cp := &MbChildSubProcessParams{}
	wfi, err := gs3wfs(ctx, cp)
	t.Require().Nil(err)
	t.Require().NotNil(wfi)
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
	wfr.WorkflowOverrides.IsUsingFlows = true
	wfr.WorkflowOverrides.WorkflowRunName = "run-demo"
	wfr.Org = t.Ou

	cp := &MbChildSubProcessParams{}
	wfi, err := s3ws(ctx, cp, &wfr)
	t.Require().Nil(err)
	t.Require().NotNil(wfi)
}
