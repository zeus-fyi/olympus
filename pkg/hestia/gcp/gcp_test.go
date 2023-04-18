package hestia_gcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type GcpTestSuite struct {
	test_suites_base.TestSuite
	g GcpClient
}

var ctx = context.Background()

func (s *GcpTestSuite) SetupTest() {
	s.InitLocalConfigs()
	g, err := InitGcpClient(ctx)
	s.Require().NoError(err)
	s.g = g
	s.Require().NotNil(s.g.Service)
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	//apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
}
func (s *GcpTestSuite) TestListSizes() {

}

func TestGcpTestSuite(t *testing.T) {
	suite.Run(t, new(GcpTestSuite))
}
