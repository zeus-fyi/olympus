package create_bases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"k8s.io/apimachinery/pkg/util/rand"
)

var ctx = context.Background()

type CreateBasesTestSuite struct {
	test_suites.DatastoresTestSuite
}

const LocalEthereumBeaconClusterDefinitionID = 1670201797184939008
const UnclassifiedLocalTestSystemID = 0

func (s *CreateBasesTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(context.Background(), s.Tc.LocalDbPgconn)
}

func (s *CreateBasesTestSuite) TestInsertBaseDefinition() {
	b := bases.NewBaseClassTopologyType(s.Tc.ProductionLocalTemporalOrgID, UnclassifiedLocalTestSystemID, "base-"+rand.String(5))
	err := InsertBase(ctx, &b)
	s.Require().Nil(err)
	s.Require().NotZero(b.TopologySystemComponentID)
}

func (s *CreateBasesTestSuite) TestInsertBasesToExistingCluster() {
	b := bases.NewBaseClassTopologyType(s.Tc.ProductionLocalTemporalOrgID, UnclassifiedLocalTestSystemID, "base-"+rand.String(5))
	b2 := bases.NewBaseClassTopologyType(s.Tc.ProductionLocalTemporalOrgID, UnclassifiedLocalTestSystemID, "base-"+rand.String(5))
	basesInsert := []bases.Base{b, b2}
	err := InsertBases(ctx, s.Tc.ProductionLocalTemporalOrgID, "unclassified-cluster", basesInsert)
	s.Require().Nil(err)
}

func TestCreateBasesTestSuite(t *testing.T) {
	suite.Run(t, new(CreateBasesTestSuite))
}
