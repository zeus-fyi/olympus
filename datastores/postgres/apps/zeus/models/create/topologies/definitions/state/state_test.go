package create_topology_deployment_status

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type CreateTopologyStateTestSuite struct {
	h hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (s *CreateTopologyStateTestSuite) TestInsertTopologyState() {
	ctx := context.Background()
	topID, _ := s.SeedTopology()
	topState := topology_deployment_status.NewTopologyStatus()
	topState.TopologyID = topID
	topState.TopologyStatus = "InProgress"
	err := InsertOrUpdateStatus(ctx, &topState)
	s.Require().Nil(err)
	s.Assert().NotEmpty(topState.UpdatedAt)

	topState.TopologyStatus = "Done"
	err = InsertOrUpdateStatus(ctx, &topState)
	s.Require().Nil(err)
	s.Assert().NotEmpty(topState.UpdatedAt)
}

func TestCreateTopologyStateTestSuite(t *testing.T) {
	suite.Run(t, new(CreateTopologyStateTestSuite))
}
