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

func (s *OrchestrationsTestSuite) TestInsertWorkflowWithComponents() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	newTemplate := WorkflowTemplate{
		WorkflowName:              "Example Workflow2",
		FundamentalPeriod:         5,
		WorkflowGroup:             "TestGroup2",
		FundamentalPeriodTimeUnit: "days",
	}

	//wt := WorkflowTasks{
	//	AggTasks: []AggTask{
	//		{
	//			AggId:      1701657830780669952,
	//			CycleCount: 1,
	//			Tasks: []AITaskLibrary{
	//				{
	//					TaskID:     1701657822027992064,
	//					CycleCount: 1,
	//					RetrievalDependencies: []RetrievalItem{
	//						{
	//							RetrievalID: 1701653245709972992,
	//						},
	//					},
	//				},
	//				{
	//					TaskID:     1701657795016150016,
	//					CycleCount: 2,
	//					RetrievalDependencies: []RetrievalItem{
	//						{
	//							RetrievalID: 1701667784112279040,
	//						},
	//					},
	//				},
	//			},
	//		},
	//	},
	//	AnalysisOnlyTasks: []AITaskLibrary{
	//		{
	//			TaskID:     1701657822027992064,
	//			CycleCount: 1,
	//			RetrievalDependencies: []RetrievalItem{
	//				{
	//					RetrievalID: 1701667813254964224,
	//				},
	//			},
	//		},
	//	},
	//}

	wt := WorkflowTasks{
		AggTasks: []AggTask{
			{
				AggId:      1701657830780669952,
				CycleCount: 1,
				Tasks: []AITaskLibrary{
					{
						TaskID:     1701657822027992064,
						CycleCount: 1,
						RetrievalDependencies: []RetrievalItem{
							{
								RetrievalID: 1701653245709972992,
							},
						},
					},
					{
						TaskID:     1701657795016150016,
						CycleCount: 2,
						RetrievalDependencies: []RetrievalItem{
							{
								RetrievalID: 1701667784112279040,
							},
						},
					},
				},
			},
			{
				AggId:      1701924144891567872,
				CycleCount: 1,
				Tasks: []AITaskLibrary{
					{
						TaskID:     1701657822027992064,
						CycleCount: 1,
						RetrievalDependencies: []RetrievalItem{
							{
								RetrievalID: 1701653245709972992,
							},
						},
					},
					{
						TaskID:     1701657795016150016,
						CycleCount: 2,
						RetrievalDependencies: []RetrievalItem{
							{
								RetrievalID: 1701667784112279040,
							},
						},
					},
				},
			},
		},
		AnalysisOnlyTasks: []AITaskLibrary{
			{
				TaskID:     1701657822027992064,
				CycleCount: 1,
				RetrievalDependencies: []RetrievalItem{
					{
						RetrievalID: 1701667813254964224,
					},
				},
			},
		},
	}
	err := InsertWorkflowWithComponents(ctx, ou, &newTemplate, wt)
	s.Require().Nil(err)
	s.Require().NotZero(newTemplate.WorkflowTemplateID)
}
