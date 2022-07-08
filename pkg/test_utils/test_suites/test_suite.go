package test_suites

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zeus-fyi/olympus/configs"
)

type BaseTestSuite struct {
	suite.Suite
	Tc configs.TestContainer
}

func (s *BaseTestSuite) SetupTest() {
	s.Tc = configs.InitLocalTestConfigs()
}

func (s *BaseTestSuite) SkipTest(b bool) {
	if b {
		s.T().SkipNow()
	}
}
func TestBaseTestSuite(t *testing.T) {
	suite.Run(t, new(BaseTestSuite))
}
