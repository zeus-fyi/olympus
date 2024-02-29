package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
)

func (t *ZeusWorkerTestSuite) TestRunAiChildAnalysisProcessWorkflow() {
	t.initWorker()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)

	aiSp := hera_search.AiSearchParams{
		TimeRange: "30 days",
	}
	hera_search.TimeRangeStringToWindow(&aiSp)
	sr, err := hera_search.SearchTwitter(ctx, t.Ou, aiSp)
	t.Require().Nil(err)
	t.Require().NotEmpty(sr)

	wfName := "twitter-extract-test-wf-eval"
	res1, err := artemis_orchestrations.SelectWorkflowTemplateByName(ctx, t.Ou, wfName)
	t.Require().Nil(err)
	t.Require().NotEmpty(res1)

	var wtvs []artemis_orchestrations.WorkflowTemplate
	for _, v := range res1.WorkflowTemplateSlice {
		wtvs = append(wtvs, v.WorkflowTemplate)
	}
	t.Require().NotNil(wtvs)
	resp, rerr := artemis_orchestrations.GetAiOrchestrationParams(ctx, t.Ou, &aiSp.Window, wtvs)
	t.Require().Nil(rerr)
	t.Require().NotNil(resp)
	t.Require().Greater(len(resp), 0)
	//err = ZeusAiPlatformWorker.ExecuteRunAiWorkflowProcess(ctx, t.Ou, resp[0])
	//t.Require().Nil(err)
}
