package read_topology

import (
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type TopologyTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *TopologyTestSuite) TestSelectTopology() {

}

func TestTopologyTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyTestSuite))
}
