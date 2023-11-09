package kronos_helix

func (t *KronosWorkerTestSuite) TestAiWorkflow() {
	ta := t.Tc.DevTemporalAuth
	//ns := "kronos.ngb72"
	//hp := "kronos.ngb72.tmprl.cloud:7233"
	//ta.Namespace = ns
	//ta.HostPort = hp
	InitKronosHelixWorker(ctx, ta)
	cKronos := KronosServiceWorker.Worker.ConnectTemporalClient()
	defer cKronos.Close()
	KronosServiceWorker.Worker.RegisterWorker(cKronos)
	err := KronosServiceWorker.Worker.Start()
	t.Require().Nil(err)

	err = KronosServiceWorker.ExecuteKronosWorkflow(ctx)
	t.Require().Nil(err)
}
