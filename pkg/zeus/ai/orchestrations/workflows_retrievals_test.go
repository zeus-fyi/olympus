package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
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
		Sg:          nil,
		Wr:          nil,
		RetryPolicy: nil,
	}
	err := ZeusAiPlatformWorker.ExecuteRetrievalsWorkflow(ctx, tte)
	t.Require().Nil(err)
}
