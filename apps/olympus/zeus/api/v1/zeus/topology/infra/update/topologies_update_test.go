package update_infra

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyUpdateActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyUpdateActionRequestTestSuite) TestUpdateChart() {
}

func TestTopologyUpdateActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyUpdateActionRequestTestSuite))
}
