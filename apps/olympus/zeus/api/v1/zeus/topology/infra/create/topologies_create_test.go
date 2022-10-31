package create_infra

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyCreateActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyCreateActionRequestTestSuite) TestCreateChart() {
	topologyActionRequest := base.TopologyActionRequest{
		Action: "create",
	}
	t.PostTopologyRequest(topologyActionRequest, 200)
}

func TestTopologyCreateActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyCreateActionRequestTestSuite))
}
