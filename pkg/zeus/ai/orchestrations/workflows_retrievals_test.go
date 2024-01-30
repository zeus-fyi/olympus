package ai_platform_service_orchestrations

import (
	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
)

func (t *ZeusWorkerTestSuite) TestRetrievalsWorkflow() {
	t.initWorker()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	act := NewZeusAiPlatformActivities()
	rets, err := act.SelectRetrievalTask(ctx, t.Ou, 1706487709357339000)
	t.Require().Nil(err)
	t.Require().NotEmpty(rets)
	ret := rets[0]

	ret.RetrievalPlatform = apiApproval

	//t.Require().Equal(apiApproval, ret.RetrievalPlatform)
	t.Require().NotNil(ret.WebFilters)
	t.Require().NotNil(ret.WebFilters.RoutingGroup)

	tte := TaskToExecute{
		Ou: t.Ou,
		Tc: TaskContext{
			TaskID:    0,
			EvalID:    1704066747085827000,
			Retrieval: ret,
			TriggerActionsApproval: artemis_orchestrations.TriggerActionsApproval{
				TriggerAction:    apiApproval,
				ApprovalID:       1706566091973007000,
				EvalID:           1704066747085827000,
				TriggerID:        1706487755984811000,
				WorkflowResultID: 1706421224945827000,
			},
			AIWorkflowTriggerResultApiResponse: artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse{
				ResponseID:  1706566091973014000,
				ApprovalID:  1706566091973007000,
				TriggerID:   1706487755984811000,
				RetrievalID: 1706487709357339000,
				ReqPayloads: []echo.Map{
					{
						"key1": "value1",
					},
					{
						"key2": "value2",
					},
				},
				RespPayloads: nil,
			},
		},
		Sg:          &hera_search.SearchResultGroup{},
		RetryPolicy: nil,
	}
	err = ZeusAiPlatformWorker.ExecuteRetrievalsWorkflow(ctx, tte)
	t.Require().Nil(err)
}
