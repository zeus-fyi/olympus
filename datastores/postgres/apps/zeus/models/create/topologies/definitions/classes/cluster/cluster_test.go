package create_clusters

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type CreateClustersTestSuite struct {
	h hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (s *CreateClustersTestSuite) TestInsertTopologyState() {
	ctx := context.Background()
	oid, uid := s.h.NewTestOrgAndUser()
	topID, _ := s.SeedTopology()
	c := NewCreateCluster()
	c.SetTopologyID(topID)
	c.SetOrgUserIDs(oid, uid)
	q := sql_query_templates.NewQueryParam("InsertCluster", "many", "where", 1000, []string{})

	// TODO
	err := c.InsertCluster(ctx, q)

	s.Require().Nil(err)
}

func TestCreateClustersTestSuite(t *testing.T) {
	suite.Run(t, new(CreateClustersTestSuite))
}
