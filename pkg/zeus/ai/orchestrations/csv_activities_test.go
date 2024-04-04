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
	// create standard processing format
	/*
		FanOutApiCallRequestTask
		---start
			wio, werr := gs3wfs(ctx, cp)
			if werr != nil {
				log.Err(werr).Msg("TokenOverflowReduction: failed to select workflow io")
				return nil, werr
			}
		---end
			_, err := s3ws(ctx, cp, wio)
			if err != nil {
				log.Err(err).Msg("TokenOverflowReduction: failed to update workflow io")
				return nil, err
			}
	*/
	// refactor JsonOutputTaskWorkflow - CreateJsonOutputModelResponse| SaveTaskOutput |
	// refactor AiAggregateAnalysisRetrievalTask -agg input
	// refactor TokenOverflowReduction -agg
	// refactor SaveTaskOutput
}
