package iris_models

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris/models/bases/autogen"
)

func (s *IrisTestSuite) TestInsertOrgRoute() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	for _, u := range s.Tc.QuikNodeURLS.Routes {
		or := iris_autogen_bases.OrgRoutes{
			RouteID:   ts.UnixTimeStampNow(),
			OrgID:     s.Tc.ProductionLocalTemporalOrgID,
			RoutePath: u,
		}
		err := InsertOrgRoute(ctx, or)
		s.Require().Nil(err)
	}
}

func (s *IrisTestSuite) TestDeleteOrgGroupRoutes() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	routes := []string{"https://chaotic-polished-glitter.quiknode.pro/2ebe51720d5565be36f36296600a571c67d60d49/",
		"https://proud-convincing-sheet.quiknode.pro/04428c76b26b5dd1be7808bd5c0df8d8dd25e86f/"}

	err := DeleteOrgRoutesFromGroup(ctx, s.Tc.ProductionLocalTemporalOrgID, "test", routes)
	s.Require().Nil(err)
}

func (s *IrisTestSuite) TestInsertOrgRoutes() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	r1 := "https://test.com/v1"
	r2 := "https://test.com/v2"
	routes := []iris_autogen_bases.OrgRoutes{
		{
			RouteID:   ts.UnixTimeStampNow(),
			RoutePath: r1,
		},
		{
			RouteID:   ts.UnixTimeStampNow(),
			RoutePath: r2,
		},
	}
	err := InsertOrgRoutes(ctx, s.Tc.ProductionLocalTemporalOrgID, routes)
	s.Require().Nil(err)

	selectRoutes, err := SelectOrgRoutes(ctx, s.Tc.ProductionLocalTemporalOrgID)
	s.Require().Nil(err)
	s.Require().NotNil(routes)

	count := 0
	for _, r := range selectRoutes {
		if r.RoutePath == r1 {
			count += 1
		}
		if r.RoutePath == r2 {
			count += 10
		}
	}
	s.Require().Equal(11, count)
	ogr := iris_autogen_bases.OrgRouteGroups{
		OrgID:          s.Tc.ProductionLocalTemporalOrgID,
		RouteGroupName: "testGroup",
	}
	err = InsertOrgRouteGroup(ctx, ogr, routes)
	s.Require().Nil(err)

	groupedRoutes, err := SelectOrgRoutesByOrgAndGroupName(ctx, s.Tc.ProductionLocalTemporalOrgID, ogr.RouteGroupName)
	s.Require().Nil(err)
	s.Require().NotNil(groupedRoutes)
	count = 0
	for _, rts := range groupedRoutes.Map {
		for _, rt := range rts {
			for _, rn := range rt {
				if rn.RoutePath == r1 {
					count += 1
				}
				if rn.RoutePath == r2 {
					count += 10
				}
			}
		}
	}
	s.Require().Equal(11, count)

	err = DeleteOrgRoutes(ctx, s.Tc.ProductionLocalTemporalOrgID, []string{r1, r2})
	s.Require().Nil(err)

	latestRoutes, err := SelectOrgRoutes(ctx, s.Tc.ProductionLocalTemporalOrgID)
	s.Require().Nil(err)
	for _, lr := range latestRoutes {
		if lr.RoutePath == r1 || lr.RoutePath == r2 {
			s.Fail("route not deleted")
		}
	}

	ogr = iris_autogen_bases.OrgRouteGroups{
		OrgID:          s.Tc.ProductionLocalTemporalOrgID,
		RouteGroupName: "testGroup2",
	}
	err = InsertOrgRoutes(ctx, s.Tc.ProductionLocalTemporalOrgID, routes)
	s.Require().Nil(err)

	err = InsertOrgRouteGroup(ctx, ogr, routes)
	s.Require().Nil(err)

	err = DeleteOrgGroupAndRoutes(ctx, s.Tc.ProductionLocalTemporalOrgID, ogr.RouteGroupName)
	s.Require().Nil(err)
}

func (s *IrisTestSuite) TestSelectOrgEndpointsAndGroupTablesCount() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	resp, err := OrgEndpointsAndGroupTablesCount(ctx, s.Tc.ProductionLocalTemporalOrgID)
	s.Require().Nil(err)
	s.Require().NotNil(resp)

	fmt.Println(resp)
	s.Require().NotZero(resp.TableCount)
	s.Require().NotZero(resp.EndpointCount)
}

func (s *IrisTestSuite) TestInsertOrgRouteQuiknode() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	groupID := 100
	for _, u := range s.Tc.QuikNodeURLS.Routes {
		or := iris_autogen_bases.OrgRoutes{
			RouteID:   ts.UnixTimeStampNow(),
			OrgID:     s.Tc.ProductionLocalTemporalOrgID,
			RoutePath: u,
		}
		err := InsertOrgRoute(ctx, or)
		s.Require().Nil(err)
		group := iris_autogen_bases.OrgRoutesGroups{
			RouteGroupID: groupID,
			RouteID:      or.RouteID,
		}
		err = InsertOrgRoutesGroups(ctx, group)
		s.Require().Nil(err)
	}
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
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

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

	ogr, err := SelectAllEndpointsAndOrgGroupRoutesByOrg(ctx, s.Tc.ProductionLocalTemporalOrgID)
	s.Require().Nil(err)
	s.Require().NotNil(ogr)

	for _, r := range ogr.Map {
		fmt.Println(r)
	}
	for _, r := range ogr.Routes {
		fmt.Println(r)
	}
}
