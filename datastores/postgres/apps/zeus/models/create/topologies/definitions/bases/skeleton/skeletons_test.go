package skeleton

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type SkeletonsTestSuite struct {
	test_suites.DatastoresTestSuite
}

func (s *SkeletonsTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(context.Background(), s.Tc.LocalDbPgconn)
}

func (s *SkeletonsTestSuite) TestInsertSkeletonDefinition() {

}

func TestSkeletonsTestSuite(t *testing.T) {
	suite.Run(t, new(SkeletonsTestSuite))
}
