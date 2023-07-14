package iris_models

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

var ts chronos.Chronos

func (s *IrisTestSuite) TestInsertOrgRoute() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	or := iris_autogen_bases.OrgRoutes{
		RouteID:   ts.UnixTimeStampNow(),
		OrgID:     s.Tc.ProductionLocalTemporalOrgID,
		RoutePath: "https://zeus.fyi/iris/test/route",
	}
	err := InsertOrgRoute(ctx, or)
	s.Require().Nil(err)
}
