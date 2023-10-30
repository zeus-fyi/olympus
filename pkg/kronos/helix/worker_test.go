package kronos_helix

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type KronosWorkerTestSuite struct {
	test_suites_base.TestSuite
}

func (t *KronosWorkerTestSuite) SetupTest() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
}

func (t *KronosWorkerTestSuite) TestKronosHelixPattern() {
	ta := t.Tc.DevTemporalAuth
	ns := "kronos.ngb72"
	hp := "kronos.ngb72.tmprl.cloud:7233"
	ta.Namespace = ns
	ta.HostPort = hp
	InitKronosHelixWorker(ctx, ta)
	cKronos := KronosServiceWorker.Worker.ConnectTemporalClient()
	defer cKronos.Close()
	KronosServiceWorker.Worker.RegisterWorker(cKronos)
	err := KronosServiceWorker.Worker.Start()
	t.Require().Nil(err)

	err = KronosServiceWorker.ExecuteKronosWorkflow(ctx)
	t.Require().Nil(err)
}

func TestKronosWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(KronosWorkerTestSuite))
}
