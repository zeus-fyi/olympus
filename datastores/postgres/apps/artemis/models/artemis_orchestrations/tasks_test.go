package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
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
