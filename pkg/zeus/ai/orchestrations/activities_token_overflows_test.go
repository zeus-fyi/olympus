package ai_platform_service_orchestrations

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (t *ZeusWorkerTestSuite) TestSearchResultsTokenOverflowReduction() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID

	aiSp := hera_search.AiSearchParams{
		TimeRange: "30 days",
		Window:    artemis_orchestrations.Window{},
	}
	hera_search.TimeRangeStringToWindow(&aiSp)
	t.Require().NotEmpty(aiSp.Window)

	sr, err := hera_search.SearchTwitter(ctx, t.Ou, aiSp)
	t.Require().Nil(err)
	t.Require().NotEmpty(sr)
	fmt.Println("srLen", len(sr))
	res := hera_search.FormatSearchResultsV3(sr)
	t.Require().NotEmpty(res)
	fmt.Println("res", len(res))

	pr := &PromptReduction{
		MarginBuffer:          0.5,
		TokenOverflowStrategy: OverflowStrategyDeduce,
		PromptReductionSearchResults: &PromptReductionSearchResults{
			InPromptBody: "This is a test prompt body",
			InSearchGroup: &hera_search.SearchResultGroup{
				PlatformName:        twitterPlatform,
				ExtractionPromptExt: "",
				Model:               Gpt4JsonModel,
				ResponseFormat:      socialMediaExtractionResponseFormat,
				SearchResults:       sr,
				Window:              aiSp.Window,
			},
			OutSearchGroups: []*hera_search.SearchResultGroup{},
		},
	}
	//err = ChunkSearchResults(ctx, pr)
	//t.Require().NoError(err)
	//
	//sgOut := pr.PromptReductionSearchResults.OutSearchGroups
	//fmt.Println("sgOut", len(sgOut))
	//t.Require().NotEmpty(sgOut)
	////act := NewZeusAiPlatformActivities()}

	err = TruncateSearchResults(ctx, pr)
	t.Require().NoError(err)
	sgOut := pr.PromptReductionSearchResults.OutSearchGroups
	fmt.Println("sgOut", len(sgOut))
	t.Require().NotEmpty(sgOut)
}
