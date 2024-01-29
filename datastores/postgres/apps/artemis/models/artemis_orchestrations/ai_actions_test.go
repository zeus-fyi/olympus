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
	qp := TriggersWorkflowQueryParams{Ou: ou, EvalID: 1702959482164376000}
	res, err := SelectTriggerActionsByOrgAndOptParams(ctx, qp)
	s.Require().Nil(err)
	s.Require().NotNil(res)

	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	qp = TriggersWorkflowQueryParams{Ou: ou, EvalID: 0}
	res2, err := SelectTriggerActionsByOrgAndOptParams(ctx, qp)
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
