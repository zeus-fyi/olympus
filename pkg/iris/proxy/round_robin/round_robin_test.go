package iris_round_robin

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type IrisRoundRobinTestSuite struct {
	test_suites_base.TestSuite
}

func (s *IrisRoundRobinTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
}

func (s *IrisRoundRobinTestSuite) TestInitRoutes() {
	//InitRoutingTables(ctx)
	for i := 0; i < 10; i++ {
		routeInfo, err := GetNextRoute(s.Tc.ProductionLocalTemporalOrgID, "test")
		s.NoError(err)
		fmt.Println(routeInfo)
	}
}

func (s *IrisRoundRobinTestSuite) TestRoundRobin() {
	// TODO IrisCache here
	//SetRouteTable(1, "test", []string{"1", "2", "3"})
	for i := 0; i < 10; i++ {
		routeInfo, err := GetNextRoute(1, "test")
		s.NoError(err)
		fmt.Println(routeInfo)
	}
}

func TestIrisRoundRobinTestSuite(t *testing.T) {
	suite.Run(t, new(IrisRoundRobinTestSuite))
}
