package web3_client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	ecdsa_signer "github.com/zeus-fyi/olympus/pkg/aegis/ecdsa"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type Web3ClientTestSuite struct {
	test_suites.EncryptionTestSuite
	GoerliWeb3  Web3Client
	MainnetWeb3 Web3Client
}

func (s *Web3ClientTestSuite) SetupTest() {
	s.InitLocalConfigs()
	pkHexString := s.Tc.LocalEcsdaTestPkey
	newAccount, err := ecdsa_signer.CreateEcdsaSignerFromPk(pkHexString)
	s.Assert().Nil(err)
	s.MainnetWeb3 = NewClientWithSigner(s.Tc.MainnetNodeUrl, newAccount)
	s.GoerliWeb3 = NewClientWithSigner(s.Tc.GoerliNodeUrl, newAccount)
}

func (s *Web3ClientTestSuite) TestWebGetBalance() {
	ctx := context.Background()
	b, err := s.GoerliWeb3.GetCurrentBalance(ctx)

	s.Require().Nil(err)
	s.Assert().NotNil(b)
	s.Assert().Greater(b.Uint64(), uint64(0))

	g, err := s.GoerliWeb3.GetCurrentBalanceGwei(ctx)
	s.Require().Nil(err)
	s.Assert().NotEqual("0", g)
}

func (s *Web3ClientTestSuite) TestWeb3ConnectMainnet() {
	ctx := context.Background()
	network, err := s.GoerliWeb3.GetNetworkName(ctx)
	s.Require().Nil(err)
	s.Assert().Equal(Goerli, network)

	network, err = s.MainnetWeb3.GetNetworkName(ctx)
	s.Require().Nil(err)
	s.Assert().Equal(Mainnet, network)
}

func TestWeb3ClientTestSuite(t *testing.T) {
	suite.Run(t, new(Web3ClientTestSuite))
}
