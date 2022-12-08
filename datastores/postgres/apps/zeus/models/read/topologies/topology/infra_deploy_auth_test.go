package read_topology

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
)

type TopologyAuthTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *TopologyTestSuite) TestDeployAuth() {
	s.InitLocalConfigs()
	ctx := context.Background()
	tr := NewInfraTopologyReader()
	tr.OrgID = 7138983863666903883

	newKns := zeus_common_types.CloudCtxNs{}
	newKns.CloudProvider = "do"
	newKns.Region = "sfo3"
	newKns.Context = "context"
	newKns.Env = "test"
	newKns.Namespace = "testnamespace"
	authed, err := tr.IsOrgCloudCtxNsAuthorized(ctx, newKns)
	s.Require().Nil(err)
	s.Assert().True(authed)
}

func TestTopologyAuthTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyAuthTestSuite))
}
