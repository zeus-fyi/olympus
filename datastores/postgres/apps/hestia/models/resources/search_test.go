package hestia_compute_resources

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/zeus/zeus/cluster_resources/nodes"
	"github.com/zeus-fyi/zeus/zeus/cluster_resources/on_demand_resources"
)

func (s *DisksTestSuite) TestNodeSearch() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	ns := nodes.NodeSearchParams{
		CloudProviderRegions: on_demand_resources.CloudProviderRegions,
		ResourceMinMax: nodes.ResourceMinMax{
			Max: nodes.ResourceAggregate{
				MonthlyPrice: 500,
			},
		},
	}
	nodeSlice, err := SearchAndSelectNodes(ctx, ns)
	s.Require().NoError(err)
	s.Require().NotEmpty(nodeSlice)

	for _, nv := range nodeSlice {
		if nv.PriceMonthly > 500 {
			s.Fail(fmt.Sprintf("PriceMonthly > 500: %f", nv.PriceMonthly))
		}
	}

	ns = nodes.NodeSearchParams{
		CloudProviderRegions: on_demand_resources.CloudProviderRegions,
		ResourceMinMax: nodes.ResourceMinMax{
			Max: nodes.ResourceAggregate{
				MonthlyPrice: 500,
			},
			Min: nodes.ResourceAggregate{
				MonthlyPrice: 300,
			},
		},
	}
	nodeSlice, err = SearchAndSelectNodes(ctx, ns)
	s.Require().NoError(err)
	s.Require().NotEmpty(nodeSlice)

	for _, nv := range nodeSlice {
		if nv.PriceMonthly < 300 {
			s.Fail(fmt.Sprintf("PriceMonthly < 300: %f", nv.PriceMonthly))
		}
	}

	ns = nodes.NodeSearchParams{
		CloudProviderRegions: on_demand_resources.CloudProviderRegions,
		ResourceMinMax: nodes.ResourceMinMax{
			Max: nodes.ResourceAggregate{
				MonthlyPrice: 500,
				MemRequests:  "50Gi",
			},
			Min: nodes.ResourceAggregate{
				MemRequests:  "20Gi",
				MonthlyPrice: 100,
			},
		},
	}
	nodeSlice, err = SearchAndSelectNodes(ctx, ns)
	s.Require().NoError(err)
	s.Require().NotEmpty(nodeSlice)

	//for _, nv := range nodeSlice {
	//	fmt.Println(nv.CloudProvider, nv.Region, nv.PriceMonthly, nv.PriceHourly, nv.Vcpus, nv.Memory, nv.Disk, nv.Description)
	//}
}

func (s *DisksTestSuite) TestNodeSearch2() {
	ns := nodes.NodeSearchParams{
		CloudProviderRegions: on_demand_resources.CloudProviderRegions,
		ResourceMinMax: nodes.ResourceMinMax{
			Max: nodes.ResourceAggregate{
				MonthlyPrice: 1000,
				MemRequests:  "50Gi",
				CpuRequests:  "12",
			},
			Min: nodes.ResourceAggregate{
				MemRequests:  "20Gi",
				MonthlyPrice: 100,
				CpuRequests:  "8",
			},
		},
	}
	nodeSlice, err := SearchAndSelectNodes(ctx, ns)
	s.Require().NoError(err)
	s.Require().NotEmpty(nodeSlice)

	for _, nv := range nodeSlice {
		fmt.Println(nv.CloudProvider, nv.Region, nv.Vcpus, nv.Memory, nv.Disk, nv.Description)
		if nv.Vcpus < 8 {
			s.Fail(fmt.Sprintf("Vcpus < 8: %f", nv.Vcpus))
		}
		if nv.Vcpus > 12 {
			s.Fail(fmt.Sprintf("Vcpus > 12: %f", nv.Vcpus))
		}
	}
}
