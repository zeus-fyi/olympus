package create_or_update_deploy

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyDeployActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyDeployActionRequestTestSuite) TestDeployChart() {
	test.Kns.Namespace = "demo"
	topologyActionRequest := base.TopologyActionRequest{
		Action: "create",
	}
	t.PostTopologyRequest(topologyActionRequest, 200)
}

func TestTopologyDeployActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyDeployActionRequestTestSuite))
}
