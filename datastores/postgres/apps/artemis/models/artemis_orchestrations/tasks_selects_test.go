package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestSelectTask() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	tasks, err := SelectTask(ctx, ou, 1705949866538066000)
	s.Require().Nil(err)
	s.Require().NotEmpty(tasks)
	s.Require().NotNil(tasks[0].Schemas)
}

func (s *OrchestrationsTestSuite) TestSelectTasks() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	tasks, err := SelectTasks(ctx, ou)
	s.Require().Nil(err)
	s.Require().NotEmpty(tasks)
}
