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
		EvalID:           1703922045959259000,
		TriggerID:        1704854895051966000,
		WorkflowResultID: 1702105254543321000,
		ApprovalState:    "pending",
	}

	// Call the function to test
	err := CreateOrUpdateTriggerActionApproval(ctx, ou, &tap)
	s.Require().Nil(err)
}
