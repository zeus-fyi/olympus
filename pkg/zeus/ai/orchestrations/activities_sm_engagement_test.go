package ai_platform_service_orchestrations

import (
	"fmt"

	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
)

func (t *ZeusWorkerTestSuite) TestJsonAggJoins() {
	t.initWorker()
	aiSp := hera_search.AiSearchParams{
		TimeRange: "30 days",
	}
	hera_search.TimeRangeStringToWindow(&aiSp)
	sr, err := hera_search.SearchTwitter(ctx, t.Ou, aiSp)
	t.Require().Nil(err)
	t.Require().NotEmpty(sr)

	if len(sr) > 20 {
		sr = sr[:20]
	}
	msgMap := make(map[int]bool)
	for _, v := range sr {
		msgMap[v.UnixTimestamp] = true
	}
	fmt.Println("srLen", len(sr))
	res := hera_search.FormatSearchResultsV3(sr)
	t.Require().NotEmpty(res)

	za := NewZeusAiPlatformActivities()
	td, err := za.SelectTaskDefinition(ctx, t.Ou, 1705180617788915000)
	t.Require().Nil(err)
	t.Require().NotNil(td)
	t.Require().Equal(1, len(td))
	tv := td[0]
	t.Require().NotNil(tv)
	t.Require().NotNil(tv.Schemas)

	fd := artemis_orchestrations.ConvertToFuncDef(socialMediaEngagementResponseFormat, tv.Schemas)
	t.Require().NotNil(fd)
	t.Require().NotNil(fd.Name)
	t.Require().NotNil(fd.Parameters)

	// make(map[string]jsonschema.Definition)
	fdv, ok := fd.Parameters.(jsonschema.Definition)
	t.Require().True(ok)
	t.Require().NotNil(fdv)

	model := Gpt4JsonModel
	pr := &PromptReduction{
		MarginBuffer:          0.5,
		TokenOverflowStrategy: OverflowStrategyTruncate,
		PromptReductionSearchResults: &PromptReductionSearchResults{
			InSearchGroup: &hera_search.SearchResultGroup{
				PlatformName:       twitterPlatform,
				Model:              model,
				ResponseFormat:     socialMediaEngagementResponseFormat,
				SearchResults:      sr,
				Window:             aiSp.Window,
				FunctionDefinition: fd,
			},
		},
	}
	fmt.Println("pr", *pr)
	err = TruncateSearchResults(ctx, pr)
	t.Require().NoError(err)
	resp, err := za.AnalyzeEngagementTweets(ctx, t.Ou, pr.PromptReductionSearchResults.InSearchGroup)
	t.Require().Nil(err)
	t.Require().NotNil(resp)
}
