package authorized_clusters

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type InsertExtClusterConfigsTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

var ctx = context.Background()

func (s *InsertExtClusterConfigsTestSuite) TestInsertExtClusterConfigs() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	pyl := []K8sClusterConfig{
		{
			ExtConfigStrID: "1707021989652474001",
			ExtConfigID:    1707021989652474001,
			//CloudProvider:  "do",
			//Region:         "nyc-3",
			//Context:        "context",
			//ContextAlias:   "alias",
			//Env:            "test123",
		},
	}
	err := InsertOrUpdateK8sClusterConfigs(ctx, s.Ou, pyl)
	s.Require().Nil(err)

	pylSelects, err := SelectAuthedAndPublicClusterConfigsByOrgID(ctx, s.Ou)
	s.Require().Nil(err)
	s.Require().Len(pylSelects, 2)
}

func (s *InsertExtClusterConfigsTestSuite) TestInsertExtClusterConfigs2() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	// SelectAuthedClusterByRouteAndOrgID
	pylSelects, err := SelectAuthedAndPublicClusterConfigsByOrgID(ctx, s.Ou)
	s.Require().Nil(err)
	s.Require().Len(pylSelects, 2)
}

func (s *InsertExtClusterConfigsTestSuite) TestSelectAuthedAndPublicClusterConfigsByOrgID() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	// SelectAuthedClusterByRouteAndOrgID
	pylSelects, err := SelectAuthedAndPublicClusterConfigsByOrgID(ctx, s.Ou)
	s.Require().Nil(err)
	s.Require().Len(pylSelects, 2)
}

func TestInsertExtClusterConfigsTestSuite(t *testing.T) {
	suite.Run(t, new(InsertExtClusterConfigsTestSuite))
}
