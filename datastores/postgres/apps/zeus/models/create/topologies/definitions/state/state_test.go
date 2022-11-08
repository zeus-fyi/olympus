package create_topology_deployment_status

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type CreateTopologyStateTestSuite struct {
	h hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (s *CreateTopologyStateTestSuite) TestInsertTopologyState() {
	ctx := context.Background()
	topID, _ := s.SeedTopology()
	topState := NewCreateState()
	topState.TopologyID = topID
	topState.TopologyStatus = "InProgress"
	err := topState.InsertStatus(ctx)
	s.Require().Nil(err)
	s.Assert().NotEmpty(topState.UpdatedAt)
}

func TestCreateTopologyStateTestSuite(t *testing.T) {
	suite.Run(t, new(CreateTopologyStateTestSuite))
}
