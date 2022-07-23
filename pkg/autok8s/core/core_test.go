package autok8s_core

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type K8TestSuite struct {
	suite.Suite
	k K8Util
}

func (s *K8TestSuite) SetupTest() {
	s.k = K8Util{}
	s.k.PrintOn = true
	s.k.ConnectToK8s()
}

func (s *K8TestSuite) TestK8Namespaces() {
	nsl, err := s.k.GetNamespaces()
	s.Nil(err)
	s.Greater(len(nsl.Items), 0)
}

func (s *K8TestSuite) TestK8Contexts() {
	kctx, err := s.k.GetContexts()
	s.Nil(err)
	s.Greater(len(kctx), 0)
}

func (s *K8TestSuite) TestK8GetPVC() {
}

func (s *K8TestSuite) TestK8UpdatePVC() {
}

func (s *K8TestSuite) TestK8GetStatefulSetList() {
}

func (s *K8TestSuite) TestK8GetStatefulSet() {
}

func (s *K8TestSuite) TestK8DeleteStatefulSet() {
}

func (s *K8TestSuite) TestK8CreateStatefulSet() {
}

func (s *K8TestSuite) TestK8UpdateStatefulSet() {
}

func (s *K8TestSuite) TestK8ReadStorageClassJsonFile() {
}

func (s *K8TestSuite) TestK8CreateStorageClass() {
}

func (s *K8TestSuite) TestK8ListStorageClasses() {
}

func (s *K8TestSuite) TestK8GetStorageClass() {
}

func TestK8sTestSuiteTest(t *testing.T) {
	suite.Run(t, new(K8TestSuite))
}
