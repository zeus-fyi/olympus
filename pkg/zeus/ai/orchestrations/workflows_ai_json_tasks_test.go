package ai_platform_service_orchestrations

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (t *ZeusWorkerTestSuite) TestJsonOutputTaskWorkflow() {
	t.initWorker()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)

	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID
	act := NewZeusAiPlatformActivities()
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
	taskID := 1705819235575890000
	td, err := act.SelectTaskDefinition(ctx, ou, taskID)
	t.Require().Nil(err)
	t.Require().NotEmpty(td)
	t.Require().Greater(len(td), 0)
	tv := td[0]
	tte := TaskToExecute{
		Ou: ou,
		Tc: TaskContext{
			TaskName: tv.TaskName,
			TaskType: tv.TaskType,
			Model:    tv.Model,
			TaskID:   tv.TaskID,
		},
		Wft: artemis_orchestrations.WorkflowTemplateData{},
		Sg: &hera_search.SearchResultGroup{
			PlatformName:   twitterPlatform,
			SourceTaskID:   tv.TaskID,
			Model:          tv.Model,
			ResponseFormat: jsonFormat,
			SearchResults:  sr,
			Window:         aiSp.Window,
		},
		Wr: &artemis_orchestrations.AIWorkflowAnalysisResult{
			OrchestrationsID:      1692062857720240000,
			ResponseID:            0,
			SourceTaskID:          taskID,
			SearchWindowUnixStart: aiSp.Window.UnixStartTime,
			SearchWindowUnixEnd:   aiSp.Window.UnixEndTime,
		},
	}
	pr := &PromptReduction{
		MarginBuffer:          0.5,
		TokenOverflowStrategy: OverflowStrategyTruncate,
		PromptReductionSearchResults: &PromptReductionSearchResults{
			InSearchGroup: tte.Sg,
		},
	}
	err = TruncateSearchResults(ctx, pr)
	t.Require().NoError(err)
	t.Require().NotEmpty(pr.PromptReductionSearchResults.OutSearchGroups)
	t.Require().NotEmpty(pr.PromptReductionSearchResults.OutSearchGroups[0].SearchResults)
	t.Require().NotEmpty(pr.PromptReductionSearchResults.OutSearchGroups[0].SearchResultChunkTokenEstimate)
	sg := pr.PromptReductionSearchResults.OutSearchGroups[0]
	tte.Sg = sg

	resp, err := ZeusAiPlatformWorker.ExecuteJsonOutputTaskWorkflow(ctx, tte)
	t.Require().Nil(err)
	t.Require().NotNil(resp)
	t.Assert().NotZero(resp.Response.ID)
	t.Assert().NotEmpty(resp.Response)
	t.Assert().NotEmpty(resp.JsonResponseResults)

	for _, res := range resp.JsonResponseResults {
		for _, v := range res {
			t.Require().NotNil(v)

			for _, f := range v.Fields {
				switch f.FieldName {
				case "msg_ids":
					t.Require().NotNil(f.IntValueSlice)
					fmt.Println("msg_ids", f.IntValueSlice)
				case "score":
					t.Require().NotNil(f.IntValue)
					fmt.Println("score", *f.IntValue)
					t.Require()
				}
			}
		}
		t.Require().NotNil(res)
	}
}
