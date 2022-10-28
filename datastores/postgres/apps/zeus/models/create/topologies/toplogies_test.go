package create_topologies

import (
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type TopologiesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (p *TopologiesTestSuite) TestInsert() {
}

func TestTopologiesTestSuite(t *testing.T) {
	suite.Run(t, new(TopologiesTestSuite))
}
