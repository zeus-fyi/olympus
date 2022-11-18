package api_test_suites

import (
	"testing"

	"github.com/stretchr/testify/suite"
	test_base "github.com/zeus-fyi/olympus/test"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

type BaseTestSuite struct {
	suite.Suite
}

func (s *BaseTestSuite) ChangeToTestDir() {
	test_base.ForceDirToTestDirLocation()
}

func (s *BaseTestSuite) TestConfigReader() {
	tc := api_configs.InitLocalTestConfigs()
	s.Assert().Equal("local", tc.Env)
}

func TestBaseTestSuite(t *testing.T) {
	suite.Run(t, new(BaseTestSuite))
}
