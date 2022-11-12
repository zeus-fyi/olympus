package web3_client

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type Web3ClientTestSuite struct {
	test_suites.PGTestSuite
}

func (s *Web3ClientTestSuite) TestWeb3Connect() {
	nodeURL := s.Tc.LocalBeaconConn
	NewClient(nodeURL)
}

func TestWeb3ClientTestSuite(t *testing.T) {
	suite.Run(t, new(Web3ClientTestSuite))
}
