package zeus_core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type ServicesTestSuite struct {
	K8TestSuite
}

func (s *ServicesTestSuite) TestGetServices() {
	ctx := context.Background()
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "dev-sfo3-zeus", Namespace: "eth-indexer"}

	svc, err := s.K.GetServiceWithKns(ctx, kns, "svc", nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(svc)
}

func TestServicesTestSuite(t *testing.T) {
	suite.Run(t, new(ServicesTestSuite))
}
