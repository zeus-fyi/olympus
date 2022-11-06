package delete_deploy

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyDeleteWorkloadActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyDeleteWorkloadActionRequestTestSuite) TestDeleteWorkloadChart() {
	test.Kns.Namespace = "demo"
	topologyActionRequest := base.TopologyActionRequest{
		Action: "delete",
	}
	t.PostTopologyRequest(topologyActionRequest, 200)
}

func TestTopologyDeleteWorkloadActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyDeleteWorkloadActionRequestTestSuite))
}
