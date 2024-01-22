package artemis_orchestrations

import (
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

	res2, err := SelectTriggerActionsByOrgAndOptParams(ctx, ou, 0)
	s.Require().Nil(err)
	s.Require().NotNil(res2)
}

func (s *OrchestrationsTestSuite) TestTriggerActions() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	// Create a TriggerActions instance
	tap := TriggerActionsApproval{
		EvalID:           1703922045959259000,
		TriggerID:        1704854895051966000,
		WorkflowResultID: 1702105254543321000,
		ApprovalState:    "pending",
	}

	// Call the function to test
	err := CreateOrUpdateTriggerActionApproval(ctx, ou, &tap)
	s.Require().Nil(err)
}

func (s *OrchestrationsTestSuite) TestSelectTriggerActions() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	res, err := SelectTriggerActionsByOrgAndOptParams(ctx, ou, 1702959482164376000)
	s.Require().Nil(err)
	s.Require().NotNil(res)

	res2, err := SelectTriggerActionsByOrgAndOptParams(ctx, ou, 0)
	s.Require().Nil(err)
	s.Require().NotNil(res2)
}

func (s *OrchestrationsTestSuite) TestSelectActionApprovals() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	res, err := SelectTriggerActionApprovals(ctx, ou, "pending")
	s.Require().Nil(err)
	s.Require().NotNil(res)
}