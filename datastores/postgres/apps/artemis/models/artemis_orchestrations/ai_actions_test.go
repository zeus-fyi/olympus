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
	err := CreateOrUpdateAction(ctx, ou, &triggerAction)
	s.Require().Nil(err)
}
