package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestTriggerActionsApprovals() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	// Create a TriggerActions instance
	tap := TriggerActionsApproval{
		EvalID:           1704066747085827000,
		TriggerID:        1706487755984811000,
		WorkflowResultID: 1705953610460614000,
		ApprovalState:    "pending",
	}

	// Call the function to test
	err := CreateOrUpdateTriggerActionApproval(ctx, ou, tap)
	s.Require().Nil(err)
}

func (s *OrchestrationsTestSuite) TestSelectTriggerActionApproval() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	ap, err := SelectTriggerActionApproval(ctx, s.Ou, "pending", 1706559169680858000)
	s.Require().Nil(err)
	s.Require().NotNil(ap)
}
