package ai_platform_service_orchestrations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type ZeusWorkerTestSuite struct {
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (t *ZeusWorkerTestSuite) SetupTest() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
}

func (t *ZeusWorkerTestSuite) TestInitWorker() {
	ta := t.Tc.DevTemporalAuth
	InitZeusAiServicesWorker(ctx, ta)
	cKronos := ZeusAiPlatformWorker.Worker.ConnectTemporalClient()
	defer cKronos.Close()
	ZeusAiPlatformWorker.Worker.RegisterWorker(cKronos)
	err := ZeusAiPlatformWorker.Worker.Start()
	t.Require().Nil(err)
}

func TestHestiaAiWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(ZeusWorkerTestSuite))
}

func (t *ZeusWorkerTestSuite) TestAiWorkflow() {
	ta := t.Tc.DevTemporalAuth
	InitZeusAiServicesWorker(ctx, ta)
	cZ := ZeusAiPlatformWorker.Worker.ConnectTemporalClient()
	defer cZ.Close()
	ZeusAiPlatformWorker.Worker.RegisterWorker(cZ)
	err := ZeusAiPlatformWorker.Worker.Start()
	t.Require().Nil(err)

	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = 7138958574876245565
	hermes_email_notifications.InitNewGmailServiceClients(ctx, t.Tc.GcpAuthJson)
	msgs, err := hermes_email_notifications.AIEmailUser.GetReadEmails("ai@zeus.fyi", 10)
	t.Require().Nil(err)

	err = ZeusAiPlatformWorker.ExecuteAiTaskWorkflow(ctx, ou, msgs)
	t.Require().Nil(err)
}
