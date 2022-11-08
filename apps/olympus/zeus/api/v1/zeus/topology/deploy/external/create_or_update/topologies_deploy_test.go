package create_or_update_deploy

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyDeployActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyDeployActionRequestTestSuite) TestDeployChart() {
}

func TestTopologyDeployActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyDeployActionRequestTestSuite))
}
