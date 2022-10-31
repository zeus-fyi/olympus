package coreK8s

import (
	"testing"

	"github.com/stretchr/testify/suite"
	clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/cluster"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/zeus_pkg"
)

type TopologyCreateActionRequestTestSuite struct {
	TopologyActionRequestTestSuite
}

func (t *TopologyCreateActionRequestTestSuite) TestCreateChart() {
	topologyActionRequest := TopologyActionRequest{
		Action:     "create",
		K8sRequest: zeus_pkg.K8sRequest{Kns: kns},
		Cluster:    clusters.NewCluster(),
	}
	t.postTopologyRequest(topologyActionRequest, 200)
}

func TestTopologyCreateActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyCreateActionRequestTestSuite))
}
