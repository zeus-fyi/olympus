package zeus_core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type NamespaceWorkloadTestSuite struct {
	K8TestSuite
}

func (s *NamespaceWorkloadTestSuite) TestGetNamespaceWorkload() {
	ctx := context.Background()
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "ethereum"
	w, err := s.K.GetWorkloadAtNamespace(ctx, kns)
	s.Require().Nil(err)
	s.Require().NotEmpty(w)
}

func TestNamespaceWorkloadTestSuite(t *testing.T) {
	suite.Run(t, new(NamespaceWorkloadTestSuite))
}
