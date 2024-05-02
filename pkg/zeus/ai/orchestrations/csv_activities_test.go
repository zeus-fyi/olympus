package ai_platform_service_orchestrations

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	utils_csv "github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/csv"
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
	wfName := "csv-analysis-16a7d5d4-b6e2"
	ue := artemis_entities.UserEntity{
		Platform: "csv-exports",
		Nickname: wfName,
	}
	ev, err := S3WfRunExport(ctx, t.Ou, wfName, &ue)
	t.Require().Nil(err)
	t.Require().NotEmpty(ev)

	headers := []string{
		"Company City", "Company State", "Latest Funding", "Latest Funding Amount", "Email", "Email Status", "Work Direct Phone", "# Employees",
		"SEO Description", "Mobile Phone", "City", "Total Funding", "Target AI", "First Name", "Company Name for Emails", "Email Confidence",
		"Facebook Url", "State", "Annual Revenue", "Email Sent", "Details", "Company Name Proper", "Industry", "Person Linkedin Url",
		"Company Linkedin Url", "Twitter Url", "Last Raised At", "Company", "First Phone", "Industry AI", "Compliment AI", "Last Name",
		"Home Phone", "Keywords", "Country", "Company Country", "Status", "Title", "Seniority", "Departments", "Corporate Phone", "Website",
		"Icebreaker AI",
	}

	csvS, err := utils_csv.SortCSV(*ue.MdSlice[0].TextData, headers)
	t.Require().Nil(err)
	fmt.Println(csvS)

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
