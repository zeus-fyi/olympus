package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestDeleteWorkflows() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	res, err := SelectWorkflowTemplates(ctx, ou)
	s.Require().Nil(err)
	s.Require().NotEmpty(res)
	//
	//err = DeleteWorkflowTemplates(ctx, ou, res)
	//s.Require().Nil(err)
}
