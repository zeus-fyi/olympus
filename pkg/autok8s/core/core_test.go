package autok8s_core

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CoreTestSuite struct {
	K8TestSuite
}

func (s *CoreTestSuite) SetupTest() {
	s.K = K8Util{}
	s.K.PrintOn = true
	s.K.ConnectToK8s()
}

func (s *CoreTestSuite) TestK8Namespaces() {
	nsl, err := s.K.GetNamespaces()
	s.Nil(err)
	s.Greater(len(nsl.Items), 0)
}

func (s *CoreTestSuite) TestK8Contexts() {
	kctx, err := s.K.GetContexts()
	s.Nil(err)
	s.Greater(len(kctx), 0)
}

func TestCoreTestSuite(t *testing.T) {
	suite.Run(t, new(CoreTestSuite))
}
