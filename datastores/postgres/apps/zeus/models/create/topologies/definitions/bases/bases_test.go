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
const AnyUseLocalTestSystemID = 1670201665179992064

func (s *CreateBasesTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(context.Background(), s.Tc.LocalDbPgconn)
}

func (s *CreateBasesTestSuite) TestInsertBaseDefinition() {
	b := bases.NewBaseClassTopologyType(s.Tc.ProductionLocalTemporalOrgID, AnyUseLocalTestSystemID, "base-"+rand.String(5))
	err := InsertBase(ctx, &b)
	s.Require().Nil(err)

	s.Require().NotZero(b.TopologySystemComponentID)
}

func TestCreateBasesTestSuite(t *testing.T) {
	suite.Run(t, new(CreateBasesTestSuite))
}
