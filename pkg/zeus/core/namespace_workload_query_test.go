package zeus_core

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type NamespaceWorkloadTestSuite struct {
	K8TestSuite
}

func (s *NamespaceWorkloadTestSuite) TestGetNamespaceWorkload() {
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "nginx"
	kns.CloudProvider = "ovh"
	kns.Context = zeusfyi
	w, err := s.K.GetWorkloadAtNamespace(ctx, kns)
	s.Require().Nil(err)
	s.Require().NotEmpty(w)
}

func TestNamespaceWorkloadTestSuite(t *testing.T) {
	suite.Run(t, new(NamespaceWorkloadTestSuite))
}
