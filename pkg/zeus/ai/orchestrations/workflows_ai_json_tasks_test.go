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
	taskID := 1701313525731432000
	td, err := act.SelectTaskDefinition(ctx, ou, taskID)
	t.Require().Nil(err)
	t.Require().NotEmpty(td)

	model := Gpt4JsonModel
	tte := TaskToExecute{
		Ou: ou,
		Tc: TaskContext{
			TaskName: "",
			TaskType: "",
			Model:    model,
			TaskID:   taskID,
		},
		Wft: artemis_orchestrations.WorkflowTemplateData{},
		Sg: &hera_search.SearchResultGroup{
			PlatformName:        twitterPlatform,
			SourceTaskID:        taskID,
			ExtractionPromptExt: "",
			Model:               model,
			ResponseFormat:      jsonFormat,
			SearchResults:       sr,
			Window:              aiSp.Window,
		},
		Wr: &artemis_orchestrations.AIWorkflowAnalysisResult{
			OrchestrationsID:      0, // TODO
			ResponseID:            0,
			SourceTaskID:          taskID,
			SearchWindowUnixStart: aiSp.Window.UnixStartTime,
			SearchWindowUnixEnd:   aiSp.Window.UnixEndTime,
		},
	}

	// TODO: add chunking
	resp, err := ZeusAiPlatformWorker.ExecuteJsonOutputTaskWorkflow(ctx, tte)
	t.Require().Nil(err)
	t.Require().NotNil(resp)
	t.Assert().NotEmpty(resp.Response)
	t.Assert().NotEmpty(resp.JsonResponseResults)

	for _, res := range resp.JsonResponseResults {
		for _, v := range res {
			t.Require().NotNil(v)

			for _, f := range v.Fields {
				switch f.FieldName {
				case "title":
					t.Require().NotNil(f.StringValue)
					fmt.Println("title", *f.StringValue)
				case "score":
					t.Require().NotNil(f.NumberValue)
					fmt.Println("score", *f.NumberValue)
				}
			}
		}
		t.Require().NotNil(res)
	}
}
