package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestTriggerInserts() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	// Create a TriggerAction instance
	triggerAction := TriggerAction{
		OrgID:        ou.OrgID,
		UserID:       ou.UserID,
		TriggerName:  "TestTrigger",
		TriggerGroup: "TestGroup",
		EvalTriggerActions: []EvalTriggerActions{
			{
				EvalID:               1702959482164376000, // Use an appropriate EvalID
				EvalTriggerState:     "active",
				EvalResultsTriggerOn: "conditionMet",
			},
		},
	}

	// Call the function to test
	err := CreateOrUpdateTriggerAction(ctx, ou, &triggerAction)
	s.Require().Nil(err)
}

func (s *OrchestrationsTestSuite) TestTriggerActions() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	// Create a TriggerAction instance
	tap := TriggerActionsApproval{
		EvalID:           1702959482164376000,
		TriggerID:        1704670458596054016,
		WorkflowResultID: 1701894366010212001,
		ApprovalState:    "pending",
	}

	// Call the function to test
	err := CreateOrUpdateTriggerActionApproval(ctx, tap)
	s.Require().Nil(err)
}

func (s *OrchestrationsTestSuite) TestSelectTriggerActions() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	res, err := SelectTriggerActionsByOrgAndEvalID(ctx, ou, 1702959482164376000)
	s.Require().Nil(err)
	s.Require().NotNil(res)
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
