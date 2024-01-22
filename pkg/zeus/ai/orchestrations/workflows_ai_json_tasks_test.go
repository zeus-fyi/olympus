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
	taskID := 1705949866538066000
	td, err := act.SelectTaskDefinition(ctx, ou, taskID)
	t.Require().Nil(err)
	t.Require().NotEmpty(td)
	t.Require().Greater(len(td), 0)
	tv := td[0]
	tv.ResponseFormat = socialMediaExtractionResponseFormat
	tte := TaskToExecute{
		Ou: ou,
		Tc: TaskContext{
			TaskName:       tv.TaskName,
			TaskType:       tv.TaskType,
			ResponseFormat: tv.ResponseFormat,
			Model:          tv.Model,
			TaskID:         tv.TaskID,
		},
		Wft: artemis_orchestrations.WorkflowTemplateData{},
		Sg: &hera_search.SearchResultGroup{
			PlatformName:   twitterPlatform,
			SourceTaskID:   tv.TaskID,
			Model:          tv.Model,
			ResponseFormat: tv.ResponseFormat,
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

	var jdef []*artemis_orchestrations.JsonSchemaDefinition
	for _, taskDef := range td {
		jdef = append(jdef, taskDef.Schemas...)
	}
	fd := artemis_orchestrations.ConvertToFuncDef(tte.Tc.TaskName, jdef)
	jsd := artemis_orchestrations.ConvertToJsonSchema(fd)
	t.Require().NotNil(jsd)

	for _, v := range jsd {
		for _, f := range v.Fields {
			switch f.FieldName {
			case "msg_ids":
				t.Require().Equal("array[integer]", f.DataType)
			}
		}
	}

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

const exResponse = `
{"twitter_extract_tweets":
[{"msg_ids":[1704249180,1704248033,1704248361,1704248824,1704247208],
	"msg_score":"3"},
{"msg_ids":[1704247583,1704247658,1704240059,1704246850],
	"msg_score":"1"},
{"msg_ids":[1704245645,1704245332,1704241340,1704245831,1704244479,1704243580,1704248681,
	1704248481,1704248020,1704248171,1704247702,1704247420,1704247020,1704247032,1704248301,
	1704248240,1704246896,1704246904,1704246850,1704246725,1704247000,1704246978,1704247566,
	1704247325,1704247249,1704247851,1704247581],
	"msg_score":"5"},
{"msg_ids":[1704245702,1704247702,1704247851,1704246939],
	"msg_score":"2"},
{"msg_ids":[1704248681,1704248725,1704248033,1704247672,1704247851,1704247408,1704247681,
	1704247459,1704247581,1704247032,1704247020,1704248000,1704248020,1704248824,1704247488,
	1704247256,1704247020,1704246850,1704247725,1704247736,1704247658,1704246437],
"msg_score":"4"}]}`
