package test_suites_base

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"

	"github.com/zeus-fyi/olympus/configs"
)

type TestSuite struct {
	CoreTestSuite
	Tc configs.TestContainer
	Ou org_users.OrgUser
}

func (s *TestSuite) SetupTest() {
	s.InitLocalConfigs()
}

func (s *TestSuite) InitLocalConfigs() {
	s.Tc = configs.InitLocalTestConfigs()
	ou := org_users.NewOrgUserWithID(s.Tc.ProductionLocalTemporalOrgID, s.Tc.ProductionLocalTemporalUserID)
	s.Ou = ou
}

func (s *TestSuite) InitProductionConfig() {
	s.Tc = configs.InitProductionConfigs()
}

func (s *TestSuite) InitStagingConfigs() {
	s.Tc = configs.InitStagingConfigs()
}

func (s *TestSuite) SkipTest(b bool) {
	if b {
		s.T().SkipNow()
	}
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
