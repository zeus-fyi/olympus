package zeus_core

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type K8TestSuite struct {
	base.TestSuite
	K K8Util
}

func (s *K8TestSuite) SetupTest() {
	s.K = K8Util{}
	s.K.PrintOn = true
	s.K.ConnectToK8s()
}

func TestK8sTestSuiteTest(t *testing.T) {
	suite.Run(t, new(K8TestSuite))
}
