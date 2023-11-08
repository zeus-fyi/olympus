package zeus_v1_compute_api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
	"github.com/zeus-fyi/zeus/zeus/cluster_resources/nodes"
	"github.com/zeus-fyi/zeus/zeus/cluster_resources/on_demand_resources"
)

var ctx = context.Background()

type ComputeRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *ComputeRequestTestSuite) TestSearchCompute() {
	t.Eg.POST("/search/nodes", NodeSearchHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

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
	resp, err := nodes.GetNodes(ctx, t.ZeusClient, ns)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func TestComputeRequestTestSuite(t *testing.T) {
	suite.Run(t, new(ComputeRequestTestSuite))
}
