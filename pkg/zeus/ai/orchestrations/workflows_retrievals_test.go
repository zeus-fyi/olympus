package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
)

func (t *ZeusWorkerTestSuite) TestRetrievalsWorkflow() {
	t.initWorker()
	tte := TaskToExecute{
		Ou: t.Ou,
		Ec: artemis_orchestrations.EvalContext{},
		Tc: TaskContext{
			TaskID:                             0,
			Retrieval:                          artemis_orchestrations.RetrievalItem{},
			AIWorkflowTriggerResultApiResponse: artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse{},
		},
		Wft:         artemis_orchestrations.WorkflowTemplateData{},
		Sg:          &hera_search.SearchResultGroup{},
		Wr:          nil,
		RetryPolicy: nil,
	}
	act := NewZeusAiPlatformActivities()
	rets, err := act.SelectRetrievalTask(ctx, t.Ou, 1706487709357339000)
	t.Require().Nil(err)
	t.Require().NotEmpty(rets)
	tte.Tc.Retrieval = rets[0]
	t.Require().Equal(webPlatform, tte.Tc.Retrieval.RetrievalPlatform)
	t.Require().NotNil(tte.Tc.Retrieval.WebFilters)
	t.Require().NotNil(tte.Tc.Retrieval.WebFilters.RoutingGroup)
	//err = ZeusAiPlatformWorker.ExecuteRetrievalsWorkflow(ctx, tte)
	//t.Require().Nil(err)
}
