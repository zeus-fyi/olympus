package iris_serverless

import (
	"time"

	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

func (t *IrisOrchestrationsTestSuite) TestUpdateResetTimer() {
	a := NewIrisPlatformActivities()
	err := a.ResyncServerlessRoutes(ctx, nil)
	t.Require().NoError(err)
}

var cctx = zeus_common_types.CloudCtxNs{
	CloudProvider: "ovh",
	Region:        "us-west-or-1",
	Context:       "kubernetes-admin@zeusfyi",
	Namespace:     "anvil-serverless-4d383226",
}

func (t *IrisOrchestrationsTestSuite) TestServerlessReset() {
	ta := t.Tc.DevTemporalAuth
	//ns := "production-iris.ngb72"
	//hp := "production-iris.ngb72.tmprl.cloud:7233"
	//ta.Namespace = ns
	//ta.HostPort = hp
	InitIrisPlatformServicesWorker(ctx, ta)
	ipw := IrisPlatformServicesWorker.Worker.ConnectTemporalClient()
	defer ipw.Close()
	IrisPlatformServicesWorker.Worker.RegisterWorker(ipw)
	err := IrisPlatformServicesWorker.Worker.Start()
	t.Require().Nil(err)

	err = IrisPlatformServicesWorker.ExecuteIrisServerlessPodRestartWorkflow(ctx, t.Tc.ProductionLocalTemporalOrgID, cctx, "anvil-9", AnvilServerlessRoutingTable, "sessionID", time.Minute*10)
	t.Require().Nil(err)

	err = IrisPlatformServicesWorker.EarlyStart(ctx, t.Tc.ProductionLocalTemporalOrgID, "anvil-9", AnvilServerlessRoutingTable, "sessionID")
	t.Require().Nil(err)
}
