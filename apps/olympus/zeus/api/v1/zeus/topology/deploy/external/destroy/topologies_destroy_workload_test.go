package destroy_deploy_request

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyDeleteWorkloadActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyDeleteWorkloadActionRequestTestSuite) TestDeleteWorkloadChart() {

}

func TestTopologyDeleteWorkloadActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyDeleteWorkloadActionRequestTestSuite))
}
