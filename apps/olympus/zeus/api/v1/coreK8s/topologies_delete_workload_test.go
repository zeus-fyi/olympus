package coreK8s

import (
	"testing"

	"github.com/stretchr/testify/suite"
	clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/cluster"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/zeus_pkg"
)

type TopologyDeleteWorkloadActionRequestTestSuite struct {
	TopologyActionRequestTestSuite
}

func (t *TopologyDeleteWorkloadActionRequestTestSuite) TestDeleteWorkloadChart() {
	kns.Namespace = "demo"
	topologyActionRequest := TopologyActionRequest{
		Action:     "delete-deploy",
		K8sRequest: zeus_pkg.K8sRequest{Kns: kns},
		Cluster:    clusters.NewCluster(),
	}
	t.postTopologyRequest(topologyActionRequest, 200)
}

func TestTopologyDeleteWorkloadActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyDeleteWorkloadActionRequestTestSuite))
}
