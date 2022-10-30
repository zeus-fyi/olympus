package coreK8s

import (
	"testing"

	"github.com/stretchr/testify/suite"
	clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/cluster"
)

type TopologyDeployActionRequestTestSuite struct {
	TopologyActionRequestTestSuite
}

func (t *TopologyDeployActionRequestTestSuite) TestDeployChart() {
	topologyActionRequest := TopologyActionRequest{
		Action:     "deploy",
		K8sRequest: K8sRequest{Kns: kns},
		Cluster:    clusters.NewCluster(),
	}
	t.postTopologyRequest(topologyActionRequest, 200)
}

func TestTopologyDeployActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyDeployActionRequestTestSuite))
}
