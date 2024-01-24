package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestSelectJsonSchemaByOrg() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	// get internal assignments
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	js, err := SelectJsonSchemaByOrg(ctx, ou)
	s.Require().Nil(err)
	s.Require().NotNil(js)
}
