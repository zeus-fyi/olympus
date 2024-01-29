package ai_platform_service_orchestrations

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

func (t *ZeusWorkerTestSuite) TestLookupEvalTriggerConditions() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	evalID := 1704066747085827000
	taskID := 1704001441781394000
	workflowTemplateID := 1704067291708542000
	act := NewZeusAiPlatformActivities()

	// MbChildSubProcessParams.WfExecParams.WorkflowTemplate.WorkflowTemplateID

	tq := artemis_orchestrations.TriggersWorkflowQueryParams{
		Ou:                 t.Ou,
		EvalID:             evalID,
		TaskID:             taskID,
		WorkflowTemplateID: workflowTemplateID,
	}
	ta, err := act.LookupEvalTriggerConditions(ctx, tq)
	t.Require().Nil(err)
	t.Require().NotNil(ta)

	tq.WorkflowTemplateID = 1704083131758512000
	ta, err = act.LookupEvalTriggerConditions(ctx, tq)
	t.Require().Nil(err)
	t.Require().Nil(ta)
}

func (t *ZeusWorkerTestSuite) TestCheckEvalTriggerCondition() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)

	evalID := 1704066747085827000
	act := NewZeusAiPlatformActivities()
	ta := &artemis_orchestrations.TriggerAction{
		TriggerID: 1706487755984811000,
		EvalTriggerActions: []artemis_orchestrations.EvalTriggerActions{
			{
				EvalID:               evalID,
				TriggerID:            1706487755984811000,
				EvalTriggerState:     infoState,
				EvalResultsTriggerOn: allPass,
			},
		},
	}
	emr := &artemis_orchestrations.EvalMetricsResults{
		EvalContext: artemis_orchestrations.EvalContext{
			EvalID:             1704066747085827000,
			EvalIterationCount: 0,
			AIWorkflowAnalysisResult: artemis_orchestrations.AIWorkflowAnalysisResult{
				WorkflowResultID: 1701894366010212001,
			},
			AIWorkflowEvalResultResponse: artemis_orchestrations.AIWorkflowEvalResultResponse{
				EvalID: evalID,
			},
		},
		EvalMetricsResults: []*artemis_orchestrations.EvalMetric{
			{
				EvalMetricResult: &artemis_orchestrations.EvalMetricResult{
					EvalMetricResultStrID:     nil,
					EvalMetricResultID:        nil,
					EvalResultOutcomeBool:     aws.Bool(true),
					EvalResultOutcomeStateStr: nil,
					RunningCycleNumber:        nil,
					EvalIterationCount:        nil,
					SearchWindowUnixStart:     nil,
					SearchWindowUnixEnd:       nil,
					EvalMetadata:              nil,
				},
			},
		},
	}
	tao, err := act.CheckEvalTriggerCondition(ctx, ta, emr)
	t.Require().NotNil(tao)
	t.Require().Nil(err)
	t.Require().NotNil(tao.TriggerActionsApprovals)

	emr.EvalMetricsResults[0].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(false)
	tao, err = act.CheckEvalTriggerCondition(ctx, ta, emr)
	t.Require().NotNil(tao)
	t.Require().Nil(tao.TriggerActionsApprovals)
	t.Require().Nil(err)
}
