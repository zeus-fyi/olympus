package ai_platform_service_orchestrations

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

func (t *ZeusWorkerTestSuite) TestS3GlobalOrgUpload() {
	ue, _ := t.getContactCsvMock()
	re, err := S3GlobalOrgUpload(ctx, t.Ou, &ue)
	t.Require().Nil(err)
	t.Require().NotNil(re)
}

func (t *ZeusWorkerTestSuite) testS3WfCycleStageImport() *MbChildSubProcessParams {
	ws := t.mockCsvMerge()
	cp := &MbChildSubProcessParams{
		WfExecParams: ws.WorkflowExecParams,
		Ou:           t.Ou,
		Tc: TaskContext{
			TaskName: "validate-emails",
			TaskType: "analysis",
		},
	}
	wsi, err := s3ws(ctx, cp, ws)
	t.Require().Nil(err)
	t.Require().NotNil(wsi)
	return cp
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
