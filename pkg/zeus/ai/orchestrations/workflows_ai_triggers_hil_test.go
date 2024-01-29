package ai_platform_service_orchestrations

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"

func (t *ZeusWorkerTestSuite) TestRunApprovedTriggerActionsWorkflow() {
	t.initWorker()

	approvalTaskGroup := ApprovalTaskGroup{
		WfID: "",
		Ou:   t.Ou,
		Taps: []artemis_orchestrations.TriggerActionsApproval{},
	}
	err := ZeusAiPlatformWorker.ExecuteTriggerActionsWorkflow(ctx, approvalTaskGroup)
	t.Require().Nil(err)
}
