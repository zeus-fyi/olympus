package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
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

func (t *ZeusWorkerTestSuite) TestExportWfCsv() {
	wfName := "test-wf"
	ue := artemis_entities.UserEntity{
		Platform: "csv-exports",
		Nickname: wfName,
	}
	ev, err := S3WfRunExport(ctx, t.Ou, "test-wf", &ue)
	t.Require().Nil(err)
	t.Require().NotEmpty(ev)
}

func (t *ZeusWorkerTestSuite) TestMergeCsvs() {

	// todo

	// refactor AiAggregateAnalysisRetrievalTask -agg input
	// refactor TokenOverflowReduction -agg
	// refactor SaveTaskOutput
}
