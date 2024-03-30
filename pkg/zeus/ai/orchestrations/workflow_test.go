package ai_platform_service_orchestrations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	artemis_hydra_orchestrations_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_s3"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

type ZeusWorkerTestSuite struct {
	test_suites_s3.S3TestSuite
}

var ctx = context.Background()

func (t *ZeusWorkerTestSuite) SetupTest() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	artemis_orchestration_auth.Bearer = t.Tc.ProductionLocalBearerToken
	auth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: t.Tc.AwsAccessKeySecretManager,
		SecretKey: t.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_auth.InitHydraSecretManagerAuthAWS(ctx, auth)
	//t.SetupLocalOvhS3()
}

func (t *ZeusWorkerTestSuite) initWorker() {
	ta := t.Tc.DevTemporalAuth
	//temporalAuthCfg := temporal_auth.TemporalAuth{
	//	ClientCertPath:   "/etc/ssl/certs/ca.pem",
	//	ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
	//	Namespace:        "production-zeus.ngb72",
	//	HostPort:         "production-zeus.ngb72.tmprl.cloud:7233",
	//}
	//ta.Namespace = temporalAuthCfg.Namespace
	//ta.HostPort = temporalAuthCfg.HostPort

	InitZeusAiServicesWorker(ctx, ta)
	cKronos := ZeusAiPlatformWorker.Worker.ConnectTemporalClient()
	defer cKronos.Close()
	ZeusAiPlatformWorker.Worker.RegisterWorker(cKronos)
	err := ZeusAiPlatformWorker.Worker.Start()
	t.Require().Nil(err)
}

func (t *ZeusWorkerTestSuite) TestInitWorker() {
	ta := t.Tc.DevTemporalAuth
	InitZeusAiServicesWorker(ctx, ta)
	cKronos := ZeusAiPlatformWorker.Worker.ConnectTemporalClient()
	defer cKronos.Close()
	ZeusAiPlatformWorker.Worker.RegisterWorker(cKronos)
	err := ZeusAiPlatformWorker.Worker.Start()
	t.Require().Nil(err)
	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = 7138958574876245565

	err = ZeusAiPlatformWorker.ExecuteAiRedditWorkflow(ctx, ou, "zeusfyi")
	t.Require().Nil(err)
}

func TestZeusWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(ZeusWorkerTestSuite))
}

func (t *ZeusWorkerTestSuite) TestRedditWorkflow() {
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

	err = ZeusAiPlatformWorker.ExecuteAiRedditWorkflow(ctx, ou, "zeusfyi")
	t.Require().Nil(err)
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
