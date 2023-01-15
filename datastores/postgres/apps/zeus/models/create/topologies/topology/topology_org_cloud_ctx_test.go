package create_topology

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type TopologiesOrgCloudCtxNsTestSuite struct {
	h hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (t *TopologiesOrgCloudCtxNsTestSuite) TestInsertTopologyAccessCloudCtxNs() {
	newKns := zeus_common_types.CloudCtxNs{}
	orgID := 7138983863666903883
	newKns.CloudProvider = "do"
	newKns.Region = "sfo3"
	newKns.Context = "context"
	newKns.Env = "test"
	newKns.Namespace = "testnamespace"

	nc := NewCreateTopologiesOrgCloudCtxNs(orgID, newKns)
	ctx := context.Background()
	err := nc.InsertTopologyAccessCloudCtxNs(ctx)
	t.Require().Nil(err)
}

func TestTopologiesOrgCloudCtxNsTestSuite(t *testing.T) {
	suite.Run(t, new(TopologiesOrgCloudCtxNsTestSuite))
}
