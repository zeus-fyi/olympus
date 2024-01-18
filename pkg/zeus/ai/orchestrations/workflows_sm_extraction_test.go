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

	if len(sr) > 50 {
		sr = sr[:50]
	}
	msgMap := make(map[int]bool)
	for _, v := range sr {
		msgMap[v.UnixTimestamp] = true
	}
	fmt.Println("srLen", len(sr))
	res := hera_search.FormatSearchResultsV3(sr)
	t.Require().NotEmpty(res)
	extPrompt := `extract only message ids from the following tweets if they not spam, overly commercial, and are discussion oriented`

	model := Gpt4JsonModel
	pr := &PromptReduction{
		MarginBuffer:          0.5,
		TokenOverflowStrategy: OverflowStrategyTruncate,
		PromptReductionSearchResults: &PromptReductionSearchResults{
			InSearchGroup: &hera_search.SearchResultGroup{
				PlatformName:        twitterPlatform,
				ExtractionPromptExt: extPrompt,
				Model:               model,
				ResponseFormat:      socialMediaExtractionResponseFormat,
				SearchResults:       sr,
				Window:              aiSp.Window,
			},
		},
	}
	err = TruncateSearchResults(ctx, pr)
	t.Require().NoError(err)
	t.Require().NotEmpty(pr.PromptReductionSearchResults.OutSearchGroups)
	t.Require().NotEmpty(pr.PromptReductionSearchResults.OutSearchGroups[0].SearchResults)
	t.Require().NotEmpty(pr.PromptReductionSearchResults.OutSearchGroups[0].SearchResultChunkTokenEstimate)
	sg := pr.PromptReductionSearchResults.OutSearchGroups[0]

	//na := NewZeusAiPlatformActivities()
	//cr, err := na.ExtractTweets(ctx, t.Ou, sg)
	//t.Require().Nil(err)
	//t.Require().NotNil(cr)
	//t.Require().NotEmpty(cr.FilteredMessages)
	//t.Require().NotEmpty(cr.FilteredMessages.MsgKeepIds)
	//
	//for _, msgID := range cr.FilteredMessages.MsgKeepIds {
	//	if _, ok := msgMap[msgID]; !ok {
	//		t.Fail("msgID not found in original search results", msgID, cr.FilteredMessages.MsgKeepIds)
	//	}
	//}
	//fmt.Println("kept", len(cr.FilteredMessages.MsgKeepIds), "all", len(msgMap))
	cr2, err := ZeusAiPlatformWorker.ExecuteSocialMediaExtractionWorkflow(ctx, t.Ou, sg)
	t.Require().Nil(err)
	t.Require().NotNil(cr2)
	t.Assert().NotEmpty(cr2.Response)
	t.Assert().NotEmpty(cr2.FilteredMessages)
	fmt.Println("kept", len(cr2.FilteredMessages.MsgKeepIds), "all", len(msgMap))
}
