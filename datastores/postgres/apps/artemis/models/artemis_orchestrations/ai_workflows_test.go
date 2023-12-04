package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestInsertAiWorkflow() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	newTemplate := WorkflowTemplate{
		WorkflowName:              "Example Workflow",
		FundamentalPeriod:         5,
		WorkflowGroup:             "TestGroup1",
		FundamentalPeriodTimeUnit: "days",
	}

	err := InsertWorkflowTemplate(ctx, ou, &newTemplate)
	s.Require().Nil(err)
	s.Require().NotZero(newTemplate.WorkflowTemplateID)
}

// InsertWorkflowWithComponents

func (s *OrchestrationsTestSuite) TestInsertWorkflowWithComponents() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	newTemplate := WorkflowTemplate{
		WorkflowName:              "Example Workflow",
		FundamentalPeriod:         5,
		WorkflowGroup:             "TestGroup1",
		FundamentalPeriodTimeUnit: "days",
	}
	tasks := []AITaskLibrary{
		{
			TaskID:     1701657822027992064,
			TaskType:   "analysis",
			TaskName:   "task-analysis-1",
			CycleCount: 1,
			RetrievalDependencies: []RetrievalItem{
				{
					RetrievalID: 1701653245709972992,
				},
			},
		},
		{
			TaskID:                1701657830780669952,
			TaskType:              "aggregation",
			CycleCount:            1,
			RetrievalDependencies: nil,
			TaskDependencies: []AITaskLibrary{
				{
					TaskID: 1701657822027992064,
				},
			},
		},
	}
	err := InsertWorkflowWithComponents(ctx, ou, &newTemplate, tasks)
	s.Require().Nil(err)
	s.Require().NotZero(newTemplate.WorkflowTemplateID)
}
