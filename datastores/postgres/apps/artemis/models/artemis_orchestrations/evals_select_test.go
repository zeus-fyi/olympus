package artemis_orchestrations

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestSelectEvals() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	// Example data for EvalFn and EvalMetrics
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	evs, err := SelectEvalFnsByOrgIDAndID(ctx, ou, 0)
	s.Require().Nil(err)
	s.Require().NotNil(evs)

	triggerCount := 0
	schemaCount := 0
	for _, ev := range evs {
		if ev.Schemas != nil {
			schemaCount++
		}

		if ev.TriggerActions != nil {
			triggerCount++
		}
	}
	s.Require().NotZero(schemaCount)
	fmt.Println("schemaCount", schemaCount)
	fmt.Println("triggerCount", triggerCount)
}

func (s *OrchestrationsTestSuite) TestSelectEvalByID() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	// Example data for EvalFn and EvalMetrics
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	evs, err := SelectEvalFnsByOrgIDAndID(ctx, ou, 1705173228307500000)
	s.Require().Nil(err)
	s.Require().NotNil(evs)

}
