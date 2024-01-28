package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

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
