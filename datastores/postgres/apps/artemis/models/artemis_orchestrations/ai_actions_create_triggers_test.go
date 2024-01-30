package artemis_orchestrations

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestTriggerInserts() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	// Create a TriggerActions instance
	triggerAction := TriggerAction{
		OrgID:        ou.OrgID,
		UserID:       ou.UserID,
		TriggerName:  "TestTrigger",
		TriggerGroup: "TestGroup",
		EvalTriggerActions: []EvalTriggerActions{
			{
				EvalID:               1703922045959259000, // Use an appropriate EvalID
				EvalTriggerState:     "info",
				EvalResultsTriggerOn: "all-pass",
			},
		},
	}

	// Call the function to test
	err := CreateOrUpdateTriggerAction(ctx, ou, &triggerAction)
	s.Require().Nil(err)

	qp := TriggersWorkflowQueryParams{Ou: ou}
	res2, err := SelectTriggerActionsByOrgAndOptParams(ctx, qp)
	s.Require().Nil(err)
	s.Require().NotNil(res2)
}

func (s *OrchestrationsTestSuite) TestCreateTriggerApiRetrieval() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	// Create a TriggerActions instance
	// round-robin
	triggerAction := TriggerAction{
		TriggerStrID:  "1705438732527176000",
		TriggerID:     1705438732527176000,
		OrgID:         ou.OrgID,
		UserID:        ou.UserID,
		TriggerName:   "social-media-approvals",
		TriggerGroup:  "social-media",
		TriggerAction: "social-media-engagement",
		TriggerRetrievals: []RetrievalItem{
			{
				RetrievalStrID: aws.String("1702098054592167000"),
				RetrievalID:    aws.Int(1702098054592167000),
				RetrievalName:  "test",
				RetrievalGroup: "test",
				RetrievalItemInstruction: RetrievalItemInstruction{
					RetrievalPlatform: "web",
					WebFilters: &WebFilters{
						RoutingGroup:       aws.String("routeGroupTestName"),
						LbStrategy:         aws.String("poll-table"),
						MaxRetries:         aws.Int(10),
						BackoffCoefficient: aws.Float64(2),
						EndpointRoutePath:  aws.String("/health"),
						EndpointREST:       aws.String("GET"),
					},
				},
			},
		},
		EvalTriggerAction: EvalTriggerActions{},
		EvalTriggerActions: []EvalTriggerActions{
			{
				EvalTriggerState:     "info",
				EvalResultsTriggerOn: "all-pass",
			},
		},
		TriggerActionsApprovals: nil,
	}

	// Call the function to test
	err := CreateOrUpdateTriggerAction(ctx, ou, &triggerAction)
	s.Require().Nil(err)
	qp := TriggersWorkflowQueryParams{Ou: ou}

	res2, err := SelectTriggerActionsByOrgAndOptParams(ctx, qp)
	s.Require().Nil(err)
	s.Require().NotNil(res2)
}

func (s *OrchestrationsTestSuite) TestSelectAll() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	qp := TriggersWorkflowQueryParams{Ou: s.Ou}

	res2, err := SelectTriggerActionsByOrgAndOptParams(ctx, qp)
	s.Require().Nil(err)
	s.Require().NotNil(res2)

}

func (s *OrchestrationsTestSuite) TestCreateTriggerApiRetrieval1() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	qp := TriggersWorkflowQueryParams{Ou: s.Ou}

	res2, err := SelectTriggerActionsByOrgAndOptParams(ctx, qp)
	s.Require().Nil(err)
	s.Require().NotNil(res2)
}
