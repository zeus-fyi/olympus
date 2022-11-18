package api_test_suites

import (
	"testing"

	"github.com/stretchr/testify/suite"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

type BaseTestSuiteTester struct {
	suite.Suite
}

func (s *BaseTestSuiteTester) TestConfigReader() {
	tc := api_configs.InitLocalTestConfigs()
	s.Assert().Equal("local", tc.Env)
}

func TestBaseTestSuiteTester(t *testing.T) {
	suite.Run(t, new(BaseTestSuiteTester))
}
