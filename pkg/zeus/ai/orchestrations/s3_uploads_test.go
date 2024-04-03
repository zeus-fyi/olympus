package ai_platform_service_orchestrations

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
