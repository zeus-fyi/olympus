package hestia_gcp

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
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
	fmt.Println(len(mt.Items))
	for _, m := range mt.Items {
		if NonSupported(m.Name) == 0 {
			fmt.Println("Skipping: ", m.Name)
			continue
		}
		skuLookup := m.GetSkuLookup()
		hourlyCost, monthlyCost, perr := hestia_compute_resources.SelectGcpPrices(ctx, skuLookup.Name, skuLookup.GPUType, skuLookup.GPUs, skuLookup.CPUs, skuLookup.MemGB)
		s.Require().NoError(perr)
		fmt.Println("Name", m.Name, "Hourly cost: ", hourlyCost, "Monthly cost: ", monthlyCost)
	}

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

func (s *GcpTestSuite) TestListServices() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	services, err := s.g.ListServices(ctx)
	s.Require().NoError(err)
	err = hestia_compute_resources.InsertGcpServices(ctx, services.Services)
	s.Require().NoError(err)

	for _, srv := range services.Services {
		err = s.g.ListServiceSKUs(ctx, srv.Name, srv.ServiceID)
		s.Require().NoError(err)
	}

}

func (s *GcpTestSuite) TestPriceAggregation() {
	sku := SkusResponse{
		Skus: []Sku{
			{
				PricingInfo: []PricingInfo{
					{
						PricingExpression: PricingExpression{
							UsageUnit:                "h",
							BaseUnit:                 "s",
							TieredRates:              []TierRate{{UnitPrice: Money{Nanos: 21811590}}},
							BaseUnitConversionFactor: 3600,
						},
					},
					{
						PricingExpression: PricingExpression{
							UsageUnit:                "GiBy.h",
							BaseUnit:                 "By.s",
							TieredRates:              []TierRate{{UnitPrice: Money{Nanos: 2923530}}},
							BaseUnitConversionFactor: 3865470566400,
						},
					},
				},
			},
		},
	}

	cpus := 2
	memoryGB := 8.0
	hoursPerMonth := 730.0
	expectedMonthlyCost := 48.90834
	monthlyCost, err := GetPrice(sku, cpus, memoryGB, hoursPerMonth)
	s.Require().NoError(err)
	fmt.Println(monthlyCost)
	s.Assert().True(monthlyCost >= expectedMonthlyCost*0.999 && monthlyCost <= expectedMonthlyCost*1.001)

	hourlyCost, err := GetPrice(sku, cpus, memoryGB, 1)
	s.Require().NoError(err)
	fmt.Println(hourlyCost)
}

func TestGcpTestSuite(t *testing.T) {
	suite.Run(t, new(GcpTestSuite))
}
