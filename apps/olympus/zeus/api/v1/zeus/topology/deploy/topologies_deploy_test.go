package deploy

import (
	"testing"

	"github.com/stretchr/testify/suite"
	clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/cluster"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/zeus_pkg"
)

type TopologyDeployActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyDeployActionRequestTestSuite) TestDeployChart() {
	test.Kns.Namespace = "demo"
	topologyActionRequest := base.TopologyActionRequest{
		Action:     "deploy",
		K8sRequest: zeus_pkg.K8sRequest{Kns: test.Kns},
		Cluster:    clusters.NewCluster(),
	}
	t.PostTopologyRequest(topologyActionRequest, 200)
}

func TestTopologyDeployActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyDeployActionRequestTestSuite))
}
