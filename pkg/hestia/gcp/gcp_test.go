package hestia_gcp

import (
	"context"
	"fmt"
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

func (s *GcpTestSuite) TestAddNodePool() {
	ci := GcpClusterInfo{
		ClusterName: "zeus-gcp-pilot-0",
		ProjectID:   "zeusfyi",
		Zone:        "us-central1-a",
	}
	mt := "e2-medium"
	ni := GkeNodePoolInfo{
		MachineType:      mt,
		InitialNodeCount: 1,
	}
	r, err := s.g.AddNodePool(ctx, ci, ni)
	s.Require().NoError(err)
	s.Require().NotNil(r)
}

func (s *GcpTestSuite) TestRemoveNodePool() {
	ci := GcpClusterInfo{
		ClusterName: "zeus-gcp-pilot-0",
		ProjectID:   "zeusfyi",
		Zone:        "us-central1-a",
	}
	mt := "e2-medium"
	ni := GkeNodePoolInfo{
		MachineType:      mt,
		InitialNodeCount: 1,
	}
	r, err := s.g.AddNodePool(ctx, ci, ni)
	s.Require().NoError(err)
	s.Require().NotNil(r)
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
	fmt.Println(mt)
	//
	//pf := hestia_infracost.ProductFilter{
	//	VendorName:    "gcp",
	//	Service:       "Compute Engine",
	//	ProductFamily: "Compute Instance",
	//	Region:        "us-central1",
	//	SKU:           "",
	//	AttributeFilters: []*hestia_infracost.AttributeFilter{
	//		{
	//			Key:   "machineType",
	//			Value: "n1-standard-64",
	//		},
	//	},
	//}
	//ic := hestia_infracost.InitInfraCostClient(ctx, s.Tc.InfraCostAPIKey)
	//for _, m := range mt.Items {
	//	fmt.Println(m)
	//	pf.AttributeFilters = []*hestia_infracost.AttributeFilter{
	//		{
	//			Key:   "machineType",
	//			Value: m.Name,
	//		},
	//	}
	//	err = ic.GetCost(ctx, pf)
	//	s.Require().NoError(err)
	//}
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
