package update

import (
	"testing"

	"github.com/stretchr/testify/suite"
	clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/cluster"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/core"
)

type TopologyUpdateActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyUpdateActionRequestTestSuite) TestUpdateeChart() {
	test.Kns.Namespace = "demo"

	topologyActionRequest := base.TopologyActionRequest{
		Action:     "update",
		K8sRequest: core.K8sRequest{Kns: test.Kns},
		Cluster:    clusters.NewCluster(),
	}
	t.PostTopologyRequest(topologyActionRequest, 200)
}

func TestTopologyUpdateActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyUpdateActionRequestTestSuite))
}
