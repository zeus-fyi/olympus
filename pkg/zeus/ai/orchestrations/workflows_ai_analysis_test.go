package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
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

	if len(sr) > 50 {
		sr = sr[:50]
	}
	tte := TaskToExecute{
		Ou: t.Ou,
	}
	cp, err := ZeusAiPlatformWorker.ExecuteRunAiWorkflowChildAnalysisProcess(ctx, tte)
	t.Require().Nil(err)
	t.Require().NotNil(cp)
}
