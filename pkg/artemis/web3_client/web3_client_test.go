package web3_client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
)

type Web3ClientTestSuite struct {
	test_suites_encryption.EncryptionTestSuite
	GoerliWeb3User  Web3Client
	GoerliWeb3User2 Web3Client
	MainnetWeb3User Web3Client
}

func (s *Web3ClientTestSuite) SetupTest() {
	s.InitLocalConfigs()
	pkHexString := s.Tc.LocalEcsdaTestPkey
	newAccount, err := accounts.ParsePrivateKey(pkHexString)
	s.Assert().Nil(err)

	pkHexString2 := s.Tc.LocalEcsdaTestPkey2
	secondAccount, err := accounts.ParsePrivateKey(pkHexString2)
	s.Assert().Nil(err)
	s.MainnetWeb3User = NewWeb3Client(s.Tc.MainnetNodeUrl, newAccount)

	s.GoerliWeb3User = NewWeb3Client(s.Tc.GoerliNodeUrl, newAccount)
	s.GoerliWeb3User2 = NewWeb3Client(s.Tc.GoerliNodeUrl, secondAccount)
}

func (s *Web3ClientTestSuite) TestWebGetBalance() {
	b, err := s.GoerliWeb3User.GetCurrentBalance(ctx)

	s.Require().Nil(err)
	s.Assert().NotNil(b)
	s.Assert().Greater(b.Uint64(), uint64(0))

	g, err := s.GoerliWeb3User.GetCurrentBalanceGwei(ctx)
	s.Require().Nil(err)
	s.Assert().NotEqual("0", g)

	g, err = s.GoerliWeb3User2.GetCurrentBalanceGwei(ctx)
	s.Require().Nil(err)
	s.Assert().NotEqual("0", g)
	fmt.Println(g)
}

func (s *Web3ClientTestSuite) TestWeb3ConnectMainnet() {
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
