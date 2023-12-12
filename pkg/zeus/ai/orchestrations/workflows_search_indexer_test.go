package ai_platform_service_orchestrations

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	artemis_hydra_orchestrations_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

func (t *ZeusWorkerTestSuite) TestAiLoop() {
	ta := t.Tc.DevTemporalAuth

	temporalAuthCfg := temporal_auth.TemporalAuth{
		ClientCertPath:   "/etc/ssl/certs/ca.pem",
		ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
		Namespace:        "production-zeus.ngb72",
		HostPort:         "production-zeus.ngb72.tmprl.cloud:7233",
	}
	ta.Namespace = temporalAuthCfg.Namespace
	ta.HostPort = temporalAuthCfg.HostPort

	InitZeusAiServicesWorker(ctx, ta)
	cKronos := ZeusAiPlatformWorker.Worker.ConnectTemporalClient()
	defer cKronos.Close()
	ZeusAiPlatformWorker.Worker.RegisterWorker(cKronos)
	err := ZeusAiPlatformWorker.Worker.Start()
	t.Require().Nil(err)

	ou := org_users.NewOrgUserWithID(t.Tc.ProductionLocalTemporalOrgID, t.Tc.ProductionLocalTemporalUserID)
	err = ZeusAiPlatformWorker.ExecuteAiRedditWorkflow(ctx, ou, "zeusfyi")
	t.Require().Nil(err)
}

func (t *ZeusWorkerTestSuite) TestAiSearchIndexerRedditWorkflow() {
	apps.Pg.InitPG(context.Background(), t.Tc.ProdLocalDbPgconn)
	//ou := org_users.NewOrgUserWithID(7138983863666903883, 7138958574876245567)
	artemis_orchestration_auth.Bearer = t.Tc.ProductionLocalTemporalBearerToken
	auth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: t.Tc.AwsAccessKeySecretManager,
		SecretKey: t.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_auth.InitHydraSecretManagerAuthAWS(ctx, auth)
	act := NewZeusAiPlatformActivities()
	resp, err := act.SearchRedditNewPostsUsingSubreddit(ctx, org_users.NewOrgUserWithID(t.Tc.ProductionLocalTemporalOrgID, t.Tc.ProductionLocalTemporalUserID), "mlops", nil)
	t.Require().Nil(err)
	t.Require().NotNil(resp)
}

func (t *ZeusWorkerTestSuite) TestAiSearchIndexerWorkflow() {
	apps.Pg.InitPG(context.Background(), t.Tc.ProdLocalDbPgconn)
	//ou := org_users.NewOrgUserWithID(7138983863666903883, 7138958574876245567)
	artemis_orchestration_auth.Bearer = t.Tc.ProductionLocalTemporalBearerToken
	auth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: t.Tc.AwsAccessKeySecretManager,
		SecretKey: t.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_auth.InitHydraSecretManagerAuthAWS(ctx, auth)
	act := NewZeusAiPlatformActivities()
	resp, err := act.SelectActiveSearchIndexerJobs(ctx)
	t.Require().Nil(err)
	t.Require().NotNil(resp)

	sgPlatformSeen := make(map[string]map[string]bool)
	var sisProcessed []hera_search.SearchIndexerParams
	for _, oj := range resp {
		switch oj.Platform {
		case "discord":
			if _, ok := sgPlatformSeen[oj.SearchGroupName]; !ok {
				sgPlatformSeen[oj.SearchGroupName] = make(map[string]bool)
			}
		case "reddit":
			if _, ok := sgPlatformSeen[oj.SearchGroupName]; !ok {
				sgPlatformSeen[oj.SearchGroupName] = make(map[string]bool)
			}
		case "twitter":
			if _, ok := sgPlatformSeen[oj.SearchGroupName]; !ok {
				sgPlatformSeen[oj.SearchGroupName] = make(map[string]bool)
			}
		case "telegram":
			if _, ok := sgPlatformSeen[oj.SearchGroupName]; !ok {
				sgPlatformSeen[oj.SearchGroupName] = make(map[string]bool)
			}
		}
		if _, ok := sgPlatformSeen[oj.SearchGroupName][oj.Platform]; !ok {
			sgPlatformSeen[oj.SearchGroupName][oj.Platform] = true
			sisProcessed = append(sisProcessed, oj)
		}
	}
	for _, si := range sisProcessed {
		if si.Platform != "twitter" {
			err = act.StartIndexingJob(ctx, si)
			t.Require().Nil(err)
		}
	}
}
