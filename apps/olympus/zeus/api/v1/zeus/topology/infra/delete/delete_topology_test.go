package delete_infra

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyDeleteActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyDeleteActionRequestTestSuite) TestDeleteWorkloadChart() {
}

func TestTopologyDeleteActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyDeleteActionRequestTestSuite))
}
