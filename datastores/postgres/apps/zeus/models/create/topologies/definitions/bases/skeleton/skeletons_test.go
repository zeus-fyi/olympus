package skeleton

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases/skeletons"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"k8s.io/apimachinery/pkg/util/rand"
)

type SkeletonsTestSuite struct {
	test_suites.DatastoresTestSuite
}

const (
	GenericBaseComponentsID = 1670202733617147904

	ConsensusClientsBaseComponentsID = 1670202869405165056
	ExecClientsBaseComponentsID      = 1670202869402443776
)

var ctx = context.Background()

func (s *SkeletonsTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(context.Background(), s.Tc.LocalDbPgconn)
}

func (s *SkeletonsTestSuite) TestInsertSkeletonDefinition() {
	sb := skeletons.NewSkeletonBase(s.Tc.ProductionLocalTemporalOrgID, GenericBaseComponentsID, "skeleton-base-"+rand.String(5))

	err := InsertSkeletonBase(ctx, &sb)
	s.Require().Nil(err)
	s.Require().NotZero(sb.TopologySkeletonBaseVersionID)
}

func TestSkeletonsTestSuite(t *testing.T) {
	suite.Run(t, new(SkeletonsTestSuite))
}
