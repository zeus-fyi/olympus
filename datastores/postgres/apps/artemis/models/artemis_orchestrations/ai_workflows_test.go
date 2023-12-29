package artemis_orchestrations

import (
	"github.com/aws/aws-sdk-go-v2/aws"
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
		WorkflowName:              "Test Workflow With Agg EvalFns 6",
		FundamentalPeriod:         33,
		WorkflowGroup:             "TestEvalAggFns",
		FundamentalPeriodTimeUnit: "minutes",
	}

	wt := WorkflowTasks{
		AggTasks: []AggTask{
			{
				AggId:      1701657830780669952,
				CycleCount: 8,
				EvalFns: []EvalFn{
					{
						EvalID: aws.Int(1703624059411640000),
					},
				},
				Tasks: []AITaskLibrary{
					{
						TaskID:     1701657822027992064,
						CycleCount: 4,
						RetrievalDependencies: []RetrievalItem{
							{
								RetrievalID: 1701653245709972992,
							},
						},
						EvalFns: []EvalFn{
							{
								EvalID: aws.Int(1703624059411640000),
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
						EvalFns: []EvalFn{
							{
								EvalID: aws.Int(1702959527789976000),
							},
						},
					},
				},
			},
		},
		AnalysisOnlyTasks: []AITaskLibrary{
			{
				TaskID:     1701657822027992064,
				CycleCount: 11,
				RetrievalDependencies: []RetrievalItem{
					{
						RetrievalID: 1701667813254964224,
					},
				},
			},
		},
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
	//		{
	//			AggId:      1701924144891567872,
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

	//wt := WorkflowTasks{
	//	AnalysisOnlyTasks: []AITaskLibrary{
	//		{
	//			TaskID:     1701657822027992064,
	//			CycleCount: 10,
	//			RetrievalDependencies: []RetrievalItem{
	//				{
	//					RetrievalID: 1701667813254964224,
	//				},
	//			},
	//			EvalFns: []EvalFn{
	//				{
	//					EvalID: aws.Int(1703624059411640000),
	//				},
	//			},
	//		},
	//	},
	//}

	err := InsertWorkflowWithComponents(ctx, ou, &newTemplate, wt)
	s.Require().Nil(err)
	s.Require().NotZero(newTemplate.WorkflowTemplateID)
}
