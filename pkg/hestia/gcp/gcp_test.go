package hestia_gcp

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/googleapi"
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
	mt := "e2-micro"
	ni := GkeNodePoolInfo{
		Name:             "nodepool-1699642242976434000-aed41f32",
		MachineType:      mt,
		InitialNodeCount: 1,
	}
	r, err := s.g.AddNodePool(ctx, ci, ni, nil, nil)
	s.Assert().NoError(err)

	var googErr *googleapi.Error
	ok := errors.As(err, &googErr)
	s.Require().True(ok)
	s.Require().Equal(googErr.Code, 409)
	s.Require().NotNil(r)
}

func (s *GcpTestSuite) TestAddNodePoolWithNvme() {
	ci := GcpClusterInfo{
		ClusterName: "zeus-gcp-pilot-0",
		ProjectID:   "zeusfyi",
		Zone:        "us-central1-a",
	}
	mt := "n1-highmem-16"
	ni := GkeNodePoolInfo{
		Name:             "test-node-pool-nvme-1234",
		MachineType:      mt,
		InitialNodeCount: 1,
	}
	tApp := &container.NodeTaint{
		Effect: "NO_SCHEDULE",
		Key:    "app",
		Value:  "sui-devnet-gcp",
	}
	labels := map[string]string{
		"app": "sui-devnet-gcp",
	}
	r, err := s.g.AddNodePool(ctx, ci, ni, []*container.NodeTaint{tApp}, labels)
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
		Name:             "test-node-pool-1",
		MachineType:      mt,
		InitialNodeCount: 1,
	}
	r, err := s.g.RemoveNodePool(ctx, ci, ni)
	s.Require().NoError(err)
	s.Require().NotNil(r)
}

func (s *GcpTestSuite) TestListMachineTypes() {
	//apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ci := GcpClusterInfo{
		ClusterName: "zeus-gcp-pilot-0",
		ProjectID:   "zeusfyi",
		Zone:        "us-central1-a",
	}

	mt, err := s.g.ListMachineTypes(ctx, ci, s.Tc.GcpAuthJson)
	s.Require().NoError(err)
	s.Require().NotEmpty(mt)

	var nodes autogen_bases.NodesSlice
	for _, m := range mt.Items {
		if NonSupported(m.Name) == 0 {
			fmt.Println("Skipping: ", m.Name)
			continue
		}
		skuLookup := m.GetSkuLookup()
		hourlyCost, monthlyCost, perr := hestia_compute_resources.SelectGcpPrices(ctx, skuLookup.Name, skuLookup.GPUType, skuLookup.GPUs, skuLookup.CPUs, skuLookup.MemGB)
		s.Require().NoError(perr)

		fmt.Println("Name", m.Name, "Hourly cost: ", hourlyCost, "Monthly cost: ", monthlyCost, "diskSize", m.MaximumPersistentDisksSizeGb)
		node := autogen_bases.Nodes{
			Memory:        m.MemoryMb,
			Vcpus:         skuLookup.CPUs,
			Disk:          10,
			DiskUnits:     "GB",
			PriceHourly:   hourlyCost,
			Region:        "us-central1",
			CloudProvider: "gcp",
			Description:   m.Description,
			Slug:          m.Name,
			MemoryUnits:   "MB",
			PriceMonthly:  monthlyCost,
			Gpus:          skuLookup.GPUs,
			GpuType:       skuLookup.GPUType,
		}
		nodes = append(nodes, node)
	}
	s.Assert().NotEmpty(nodes)
	err = hestia_compute_resources.InsertNodes(ctx, nodes)
	s.Require().NoError(err)
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
