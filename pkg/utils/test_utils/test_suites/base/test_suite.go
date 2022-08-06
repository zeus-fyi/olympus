package base

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zeus-fyi/olympus/configs"
)

type TestSuite struct {
	suite.Suite
	Tc configs.TestContainer
}

func (s *TestSuite) SetupTest() {
	s.InitConfigs()
}

func (s *TestSuite) InitConfigs() {
	s.Tc = configs.InitLocalTestConfigs()
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
func TestBaseTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
