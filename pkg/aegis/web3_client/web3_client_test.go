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
	GoerliWeb3User  Web3Client
	GoerliWeb3User2 Web3Client
	MainnetWeb3User Web3Client
}

func (s *Web3ClientTestSuite) SetupTest() {
	s.InitLocalConfigs()
	pkHexString := s.Tc.LocalEcsdaTestPkey
	newAccount, err := ecdsa_signer.CreateEcdsaSignerFromPk(pkHexString)
	s.Assert().Nil(err)

	pkHexString2 := s.Tc.LocalEcsdaTestPkey2
	secondAccount, err := ecdsa_signer.CreateEcdsaSignerFromPk(pkHexString2)
	s.Assert().Nil(err)
	s.MainnetWeb3User = NewClientWithSigner(s.Tc.MainnetNodeUrl, newAccount)

	s.GoerliWeb3User = NewClientWithSigner(s.Tc.GoerliNodeUrl, newAccount)
	s.GoerliWeb3User2 = NewClientWithSigner(s.Tc.GoerliNodeUrl, secondAccount)
}

func (s *Web3ClientTestSuite) TestWebGetBalance() {
	ctx := context.Background()
	b, err := s.GoerliWeb3User.GetCurrentBalance(ctx)

	s.Require().Nil(err)
	s.Assert().NotNil(b)
	s.Assert().Greater(b.Uint64(), uint64(0))

	g, err := s.GoerliWeb3User.GetCurrentBalanceGwei(ctx)
	s.Require().Nil(err)
	s.Assert().NotEqual("0", g)
}

func (s *Web3ClientTestSuite) TestWeb3ConnectMainnet() {
	ctx := context.Background()
	network, err := s.GoerliWeb3User.GetNetworkName(ctx)
	s.Require().Nil(err)
	s.Assert().Equal(Goerli, network)

	network, err = s.MainnetWeb3User.GetNetworkName(ctx)
	s.Require().Nil(err)
	s.Assert().Equal(Mainnet, network)
}

func TestWeb3ClientTestSuite(t *testing.T) {
	suite.Run(t, new(Web3ClientTestSuite))
}
