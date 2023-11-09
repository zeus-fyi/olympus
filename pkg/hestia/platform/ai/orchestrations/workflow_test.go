package ai_platform_service_orchestrations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type HestiaAiWorkerTestSuite struct {
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (t *HestiaAiWorkerTestSuite) SetupTest() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
}

func (t *HestiaAiWorkerTestSuite) TestInitWorker() {
	ta := t.Tc.DevTemporalAuth
	InitHestiaIrisPlatformServicesWorker(ctx, ta)
	cKronos := HestiaAiPlatformWorker.Worker.ConnectTemporalClient()
	defer cKronos.Close()
	HestiaAiPlatformWorker.Worker.RegisterWorker(cKronos)
	err := HestiaAiPlatformWorker.Worker.Start()
	t.Require().Nil(err)

}

func TestHestiaAiWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(HestiaAiWorkerTestSuite))
}

//func (t *KronosWorkerTestSuite) TestAiWorkflow() {
//	ta := t.Tc.DevTemporalAuth
//	//ns := "kronos.ngb72"
//	//hp := "kronos.ngb72.tmprl.cloud:7233"
//	//ta.Namespace = ns
//	//ta.HostPort = hp
//	InitKronosHelixWorker(ctx, ta)
//	cKronos := KronosServiceWorker.Worker.ConnectTemporalClient()
//	defer cKronos.Close()
//	KronosServiceWorker.Worker.RegisterWorker(cKronos)
//	err := KronosServiceWorker.Worker.Start()
//	t.Require().Nil(err)
//
//	ou := org_users.OrgUser{}
//	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
//	ou.UserID = 7138958574876245565
//	content := "write why is my golang json unmarshal not working properly only on linux"
//	err = KronosServiceWorker.ExecuteAiTaskWorkflow(ctx, ou, content)
//	t.Require().Nil(err)
//}
