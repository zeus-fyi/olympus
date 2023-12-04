package artemis_orchestrations

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

func (s *OrchestrationsTestSuite) TestInsertRetrieval() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

}
