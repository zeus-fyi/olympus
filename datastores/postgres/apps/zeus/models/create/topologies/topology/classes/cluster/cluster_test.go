package create_clusters

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	create_systems "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology/classes/systems"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"k8s.io/apimachinery/pkg/util/rand"
)

var ctx = context.Background()

type CreateClustersTestSuite struct {
	test_suites.DatastoresTestSuite
}

func (s *CreateClustersTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(context.Background(), s.Tc.LocalDbPgconn)
}

func (s *CreateClustersTestSuite) TestInsertClusterDefinition() {
	c := NewClusterClassTopologyType(s.Tc.ProductionLocalTemporalOrgID, "cluster-"+rand.String(5))
	err := create_systems.InsertSystem(ctx, &c.Systems)
	s.Require().Nil(err)
}

func TestCreateClustersTestSuite(t *testing.T) {
	suite.Run(t, new(CreateClustersTestSuite))
}
