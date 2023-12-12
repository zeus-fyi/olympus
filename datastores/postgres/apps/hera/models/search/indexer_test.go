package hera_search

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *SearchAITestSuite) TestUpdateTwitterSearchQueryStatus() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	sp := SearchIndexerParams{
		Active:          true,
		SearchGroupName: defaultTwitterSearchGroupName,
	}
	err := UpdateTwitterSearchQueryStatus(ctx, ou, sp)
	s.Require().Nil(err)
}
