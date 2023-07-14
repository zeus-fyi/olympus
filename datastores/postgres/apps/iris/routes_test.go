package iris_models

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

// todo add routes to db for beacon node lb

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

func (s *IrisTestSuite) TestInsertOrgRouteGroup() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ogr := iris_autogen_bases.OrgRouteGroups{
		RouteGroupID:   0,
		OrgID:          s.Tc.ProductionLocalTemporalOrgID,
		RouteGroupName: "test",
	}
	err := InsertOrgRouteGroup(ctx, ogr)
	s.Require().Nil(err)
}
func (s *IrisTestSuite) TestInsertOrgRoutesToGroup() {
	or := iris_autogen_bases.OrgRoutesGroups{
		RouteGroupID: 0,
		RouteID:      1689299343647752000,
	}
	err := InsertOrgRoutesGroups(ctx, or)
	s.Require().Nil(err)

	or = iris_autogen_bases.OrgRoutesGroups{
		RouteGroupID: 0,
		RouteID:      1689299347619194000,
	}
	err = InsertOrgRoutesGroups(ctx, or)
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

func (s *IrisTestSuite) TestSelectAllOrgRoutes() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	routes, err := SelectAllOrgRoutes(ctx)
	s.Require().Nil(err)
	s.Require().NotNil(routes)

	for orgID, r := range routes.Map {
		fmt.Println(orgID)
		s.Require().Equal(s.Tc.ProductionLocalTemporalOrgID, orgID)
		for k, v := range r {
			fmt.Println(k, v)
		}
	}

}
