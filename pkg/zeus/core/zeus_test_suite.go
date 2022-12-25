package zeus_core

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type K8TestSuite struct {
	test_suites_base.TestSuite
	K K8Util
}

func (s *K8TestSuite) SetupTest() {
	s.ConnectToK8s()
}

func (s *K8TestSuite) ConnectToK8s() {
	s.K = K8Util{}
	s.K.PrintOn = true
	s.K.ConnectToK8s()

	s.K.SetContext("do-nyc1-do-nyc1-zeus-demo")
}

func TestK8sTestSuiteTest(t *testing.T) {
	suite.Run(t, new(K8TestSuite))
}
