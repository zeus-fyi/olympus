package hestia_compute_resources

import (
	"github.com/zeus-fyi/zeus/zeus/cluster_resources/nodes"
	"github.com/zeus-fyi/zeus/zeus/cluster_resources/on_demand_resources"
)

func (s *DisksTestSuite) TestNodeSearchV2() {
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
	resourceMap, err := SearchAndSelectNodesV2(ctx, ns)
	s.Require().NoError(err)
	s.Require().NotEmpty(resourceMap)
}
