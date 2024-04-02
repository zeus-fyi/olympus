package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

const (
	validemailRetQp = "validemail-query-params"
)

func (t *ZeusWorkerTestSuite) TestWfCsv() {
	cp := t.testS3WfCycleStageImport()
	za := NewZeusAiPlatformActivities()
	wr := &artemis_orchestrations.AIWorkflowAnalysisResult{}
	res, err := za.SaveCsvTaskOutput(ctx, cp, wr)
	t.Require().Nil(err)
	t.Assert().NotEmpty(res)
}
