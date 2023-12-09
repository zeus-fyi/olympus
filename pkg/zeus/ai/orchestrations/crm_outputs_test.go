package ai_platform_service_orchestrations

import (
	artemis_hydra_orchestrations_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

func (t *ZeusWorkerTestSuite) TestAiTelegramTask() {
	//ou := org_users.NewOrgUserWithID(7138983863666903883, 7138958574876245567)
	artemis_orchestration_auth.Bearer = t.Tc.ProductionLocalTemporalBearerToken
	auth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: t.Tc.AwsAccessKeySecretManager,
		SecretKey: t.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_auth.InitHydraSecretManagerAuthAWS(ctx, auth)
	msgs, err := GetPandoraMessages(ctx, "Zeus")
	t.Require().Nil(err)
	t.Require().NotNil(msgs)
	//params := hera_search.AiSearchParams{
	//	GroupFilter: "Zeus <> Li",
	//}
	//	hera_openai.InitHeraOpenAI(t.Tc.OpenAIAuth)
	//	// resty_base.GetBaseRestyClient("https://api.openai.com", bearer)
	//	sr, err := hera_search.SearchTelegram(ctx, ou, params)
	//	resp, err := AiTelegramTask(ctx, ou, sr, params)
	//	t.Require().Nil(err)
	//	t.Require().NotNil(resp)
	//
}
