package ai_platform_service_orchestrations

import (
	"fmt"

	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
)

func (t *ZeusWorkerTestSuite) TestSmExtractionWfTwitter() {
	t.initWorker()
	artemis_orchestration_auth.Bearer = t.Tc.ProductionLocalBearerToken
	aiSp := hera_search.AiSearchParams{
		TimeRange: "30 days",
	}
	hera_search.TimeRangeStringToWindow(&aiSp)
	sr, err := hera_search.SearchTwitter(ctx, t.Ou, aiSp)
	t.Require().Nil(err)
	t.Require().NotEmpty(sr)

	fmt.Println("srLen", len(sr))
	res := hera_search.FormatSearchResultsV3(sr)
	t.Require().NotEmpty(res)
	fmt.Println("res", len(res))

	tc, err := GetTokenCountEstimate(ctx, Gpt4JsonModel, res)
	t.Require().Nil(err)
	t.Require().NotZero(tc)
	fmt.Println("tc", tc)
	//extPrompt := "todo, which tweets to extract"
	//sg := &hera_search.SearchResultGroup{
	//	PlatformName:        twitterPlatform,
	//	ExtractionPromptExt: extPrompt,
	//	Model:               Gpt4JsonModel,
	//	ResponseFormat:      socialMediaExtractionResponseFormat,
	//	SearchResults:       sr,
	//	Window:              aiSp.Window,
	//}
	//
	//cr, err := ZeusAiPlatformWorker.ExecuteSocialMediaExtractionWorkflow(ctx, t.Ou, sg)
	//t.Require().Nil(err)
	//t.Require().NotNil(cr)

	// TODO, verify response
}
