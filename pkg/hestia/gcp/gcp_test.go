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
	g, err := InitGcpClient(ctx, s.Tc.GcpAuthJson)
	s.Require().NoError(err)
	s.g = g
	s.Require().NotNil(s.g.Service)
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
}

func (s *GcpTestSuite) TestListMachineTypes() {
	ci := GcpClusterInfo{
		ClusterName: "zeus-gcp-pilot-0",
		ProjectID:   "zeusfyi",
		Zone:        "us-central1-a",
	}

	mt, err := s.g.ListMachineTypes(ctx, ci, s.Tc.GcpAuthJson)
	s.Require().NoError(err)
	s.Require().NotEmpty(mt)
}

func (s *GcpTestSuite) TestListNodes() {
	ci := GcpClusterInfo{
		ClusterName: "zeus-gcp-pilot-0",
		ProjectID:   "zeusfyi",
		Zone:        "us-central1-a",
	}

	ns, err := s.g.ListNodes(ctx, ci)
	s.Require().NoError(err)
	s.Require().NotNil(ns)
}

func TestGcpTestSuite(t *testing.T) {
	suite.Run(t, new(GcpTestSuite))
}
