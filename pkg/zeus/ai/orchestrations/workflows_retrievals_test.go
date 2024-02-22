package ai_platform_service_orchestrations

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
)

func (t *ZeusWorkerTestSuite) TestRetrievalsWorkflowTask() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	act := NewZeusAiPlatformActivities()
	rets, err := act.SelectRetrievalTask(ctx, t.Ou, 1706767039731058000)
	t.Require().Nil(err)
	t.Require().NotEmpty(rets)
	ret := rets[0]

	ret.RetrievalPlatform = apiApproval
	cp := &MbChildSubProcessParams{
		WfID:         uuid.New().String(),
		Ou:           t.Ou,
		WfExecParams: artemis_orchestrations.WorkflowExecParams{},
		Oj: artemis_orchestrations.OrchestrationJob{Orchestrations: artemis_autogen_bases.Orchestrations{
			OrchestrationID: 1706767039731058000,
		}},
		Window: artemis_orchestrations.Window{},
		Wsr: artemis_orchestrations.WorkflowStageReference{
			WorkflowRunID: 1704069081079680000,
			ChildWfID:     "TestRetrievalsWorkflow-" + uuid.New().String(),
		},
		Tc: TaskContext{
			TaskID:    1706842030247904000,
			Retrieval: ret,
		},
	}

	cp, err = ZeusAiPlatformWorker.ExecuteRetrievalsWorkflow(ctx, cp)
	t.Require().Nil(err)
	t.Require().NotZero(cp.Wsr.InputID)
}

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

	cp := &MbChildSubProcessParams{
		WfID:         uuid.New().String(),
		Ou:           t.Ou,
		WfExecParams: artemis_orchestrations.WorkflowExecParams{},
		Oj:           artemis_orchestrations.OrchestrationJob{},
		Window:       artemis_orchestrations.Window{},
		Wsr: artemis_orchestrations.WorkflowStageReference{
			WorkflowRunID:  0,
			ChildWfID:      "TestRetrievalsWorkflow-" + uuid.New().String(),
			RunCycle:       0,
			IterationCount: 0,
			ChunkOffset:    0,
		},
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
			},
		},
	}

	_, err = ZeusAiPlatformWorker.ExecuteRetrievalsWorkflow(ctx, cp)
	t.Require().Nil(err)
}

func (t *ZeusWorkerTestSuite) TestRetrievalsExtract() {
	m := echo.Map{
		"q":        "best book",
		"category": "fiction books",
	}

	route := "customsearch/h/v1?q={q}&category={category}"
	ps, err := ReplaceParams(route, m)
	t.Require().Nil(err)

	expected := "customsearch/h/v1?q=best+book&category=fiction+books" // Expect spaces to be replaced with '+'
	t.Require().Equal(expected, ps)
	fmt.Println(ps)
}
