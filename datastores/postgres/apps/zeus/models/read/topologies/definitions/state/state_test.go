package read_topology_deployment_status

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type ReadTopologyStateTestSuite struct {
	h hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (s *ReadTopologyStateTestSuite) TestReadLatestTopologyState() {
	apps.Pg.InitPG(context.Background(), s.Tc.ProdLocalDbPgconn)
	dr := NewReadDeploymentStatusesGroup()
	ctx := context.Background()
	orgID := 7138983863666903883
	userID := 1667452524356256466
	ou := org_users.NewOrgUserWithID(orgID, userID)

	err := dr.ReadLatestDeployedClusterTopologies(ctx, 1671248907408699000, ou)
	s.Require().Nil(err)
	s.Assert().NotEmpty(dr.Slice)
	for _, v := range dr.Slice {
		fmt.Println(v.ClusterName)
		fmt.Println(v.UpdatedAt.UTC())
	}
}

func (s *ReadTopologyStateTestSuite) TestReadTopologyState() {
	topologyID := 7155775605218483902

	dr := NewReadDeploymentStatusesGroup()
	ctx := context.Background()
	orgID := 1667452524363177528
	userID := 1667452524356256466
	ou := org_users.NewOrgUserWithID(orgID, userID)
	err := dr.ReadStatus(ctx, topologyID, ou)
	s.Require().Nil(err)
	s.Assert().NotEmpty(dr.Slice)
}

func TestCreateTopologyStateTestSuite(t *testing.T) {
	suite.Run(t, new(ReadTopologyStateTestSuite))
}
