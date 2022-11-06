package temporal_client

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type TemporalClientTestSuite struct {
	base.TestSuite
}

func (s *TemporalClientTestSuite) SetupTest() {
	s.InitLocalConfigs()
}

func (s *TemporalClientTestSuite) TestCreateClient() {
	certPath := "./zeus.fyi/ca.pem"
	pemPath := "./zeus.fyi/ca.key"
	namespace := s.Tc.DevTemporalNs
	hostPort := s.Tc.DevTemporalHostPort

	auth := TemporalAuth{
		ClientCertPath:   certPath,
		ClientPEMKeyPath: pemPath,
		Namespace:        namespace,
		HostPort:         hostPort,
	}
	err := ConnectClient(auth)
	s.Require().Nil(err)
}

func TestTemporalClientTestSuite(t *testing.T) {
	suite.Run(t, new(TemporalClientTestSuite))
}
