package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestSelectRuns() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	ojs, err := SelectAiSystemOrchestrations(ctx, ou, 0)
	s.Require().Nil(err)
	s.Assert().NotEmpty(ojs)
}

func (s *OrchestrationsTestSuite) TestSelectRun() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	ojs, err := SelectAiSystemOrchestrations(ctx, ou, 1709076571395743000)
	s.Require().Nil(err)
	s.Assert().NotEmpty(ojs)
}
