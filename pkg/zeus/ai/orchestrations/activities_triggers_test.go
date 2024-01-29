package ai_platform_service_orchestrations

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

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
				EvalTriggerState:     "info",
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
					EvalResultOutcomeBool: aws.Bool(true),
				},
			},
		},
	}
	tao, err := act.CheckEvalTriggerCondition(ctx, ta, emr)
	t.Require().NotNil(tao)
	t.Require().Nil(err)
	t.Require().NotNil(tao.TriggerActionsApprovals)
}
