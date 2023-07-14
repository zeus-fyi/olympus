package iris_models

import (
	"fmt"

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
		RoutePath: "https://zeus.fyi/iris/test/route1",
	}
	err := InsertOrgRoute(ctx, or)
	s.Require().Nil(err)
}

func (s *IrisTestSuite) TestSelectOrgRoutes() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	routes, err := SelectOrgRoutes(ctx, s.Tc.ProductionLocalTemporalOrgID)
	s.Require().Nil(err)
	s.Require().NotNil(routes)

	for _, r := range routes {
		s.Require().Equal(s.Tc.ProductionLocalTemporalOrgID, r.OrgID)
		fmt.Println(r.RoutePath)
	}
}
