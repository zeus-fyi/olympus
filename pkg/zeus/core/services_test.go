package zeus_core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServicesTestSuite struct {
	K8TestSuite
}

func (s *ServicesTestSuite) TestGetServices() {
	ctx := context.Background()
	var kns = KubeCtxNs{CloudProvider: "do", Region: "sfo3", CtxType: "dev-sfo3-zeus", Namespace: "eth-indexer"}

	svc, err := s.K.GetServiceWithKns(ctx, kns, "svc", nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(svc)
}

func TestServicesTestSuite(t *testing.T) {
	suite.Run(t, new(ServicesTestSuite))
}
