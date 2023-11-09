package kronos_helix

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"

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

	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = 7138958574876245565
	content := "write why is my golang json unmarshal not working properly only on linux"
	err = KronosServiceWorker.ExecuteAiTaskWorkflow(ctx, ou, content)
	t.Require().Nil(err)
}
