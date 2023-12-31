package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestInsertTask() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	// Test data for insertion
	testTask := AITaskLibrary{
		OrgID:                 s.Tc.ProductionLocalTemporalOrgID,
		UserID:                s.Tc.ProductionLocalTemporalUserID,
		MaxTokensPerTask:      100,
		TaskType:              "aggregation",
		TaskName:              "task-aggregation-3",
		TaskGroup:             "default",
		TokenOverflowStrategy: "deduce",
		Model:                 "gpt-4",
		Prompt:                "zzztest prompt",
	}

	err := InsertTask(ctx, &testTask)
	s.Require().Nil(err)
	s.Require().NotZero(testTask.TaskID)
}

func (s *OrchestrationsTestSuite) TestTaskAggregation() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	name := "Test Workflow With Agg EvalFns Nulls 7"
	res, err := SelectWorkflowTemplateByName(ctx, ou, name)
	s.Require().Nil(err)
	s.Require().NotEmpty(res)

	wte := ConvertTemplateValuesToWorkflowTemplateData(WorkflowTemplate{}, res.WorkflowTemplateSlice[0])
	s.Require().NotNil(wte)
	s.Require().NotEmpty(wte.CycleCountTaskRelative)
	s.Require().NotEmpty(wte.CycleCountTaskRelative.AggNormalizedCycleCounts)
	s.Assert().Equal(88, wte.CycleCountTaskRelative.AggNormalizedCycleCounts[1701313112337875000])
	s.Assert().Equal(704, wte.CycleCountTaskRelative.AggEvalNormalizedCycleCounts[1701313112337875000][1702961311357646000])
	s.Assert().Equal(22, wte.CycleCountTaskRelative.AggAnalysisEvalNormalizedCycleCounts[1701313112337875000][1701313525731432000][1702961311357646000])
}
