package web3_client

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/v4/common"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
)

type Web3ClientTestSuite struct {
	test_suites_encryption.EncryptionTestSuite
	GoerliWeb3User                Web3Client
	GoerliWeb3User2               Web3Client
	MainnetWeb3User               Web3Client
	MainnetWeb3UserExternal       Web3Client
	LocalMainnetWeb3User          Web3Client
	LocalHardhatMainnetUser       Web3Client
	HostedHardhatMainnetUser      Web3Client
	ProxyHostedHardhatMainnetUser Web3Client
}

func (s *Web3ClientTestSuite) SetupTest() {
	s.InitLocalConfigs()
	pkHexString := s.Tc.LocalEcsdaTestPkey
	newAccount, err := accounts.ParsePrivateKey(pkHexString)
	s.Assert().Nil(err)

	pkHexString2 := s.Tc.LocalEcsdaTestPkey2
	secondAccount, err := accounts.ParsePrivateKey(pkHexString2)
	s.Assert().Nil(err)
	s.MainnetWeb3UserExternal = NewWeb3Client(s.Tc.MainnetNodeUrl, newAccount)
	s.HostedHardhatMainnetUser = NewWeb3Client("https://hardhat.zeus.fyi", newAccount)
	s.GoerliWeb3User = NewWeb3Client(s.Tc.GoerliNodeUrl, newAccount)
	s.GoerliWeb3User2 = NewWeb3Client(s.Tc.GoerliNodeUrl, secondAccount)

	s.MainnetWeb3User = NewWeb3Client(s.Tc.LocalBeaconConn, newAccount)
	m := map[string]string{
		"Authorization": "Bearer " + s.Tc.ProductionLocalTemporalBearerToken,
	}
	s.MainnetWeb3User.Headers = m
	s.HostedHardhatMainnetUser.Headers = m

	s.LocalMainnetWeb3User = NewWeb3Client("http://localhost:8545", newAccount)

	newAccount, err = accounts.ParsePrivateKey("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	s.Assert().Nil(err)
	s.HostedHardhatMainnetUser.Account = newAccount
	// iris.zeus.fyi
	s.ProxyHostedHardhatMainnetUser = NewWeb3Client("https://iris.zeus.fyi/v1/internal/", newAccount)
	//s.ProxyHostedHardhatMainnetUser = NewWeb3Client("http://localhost:8080/v1/internal/", newAccount)
	s.ProxyHostedHardhatMainnetUser.Headers = map[string]string{
		"Authorization": "Bearer " + s.Tc.ProductionLocalTemporalBearerToken,
	}
	s.LocalHardhatMainnetUser.Account = newAccount
	s.LocalHardhatMainnetUser = NewWeb3Client("http://localhost:8545", newAccount)

}

func (s *Web3ClientTestSuite) TestGetProxyHardhat() {
	pb, err := s.ProxyHostedHardhatMainnetUser.GetCurrentBalance(ctx)
	s.Require().Nil(err)
	s.Assert().NotNil(pb)
	fmt.Println("bal", pb.String())

	hb, err := s.HostedHardhatMainnetUser.GetCurrentBalance(ctx)
	s.Require().Nil(err)
	s.Assert().NotNil(hb)
	fmt.Println("bal", hb.String())

	s.Assert().Equal(hb.String(), pb.String())
}

func (s *Web3ClientTestSuite) TestGetBlockHeight() {
	b, err := s.MainnetWeb3User.GetHeadBlockHeight(ctx)
	s.Require().Nil(err)
	s.Assert().NotNil(b)
	fmt.Println("blockNumber", b.String())
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

func (s *Web3ClientTestSuite) TestReadMempool() {
	s.MainnetWeb3User.Web3Actions.Dial()
	defer s.MainnetWeb3User.Close()
	mempool, err := s.MainnetWeb3User.Web3Actions.GetTxPoolContent(ctx)
	s.Require().Nil(err)
	s.Assert().NotNil(mempool)
	uswap := InitUniswapClient(ctx, s.MainnetWeb3User)
	s.Require().Nil(err)
	smartContractAddrFilter := common.HexToAddress(uswap.SmartContractAddr)
	smartContractAddrFilterString := smartContractAddrFilter.String()
	for userAddr, txPoolQueue := range mempool["pending"] {
		for order, tx := range txPoolQueue {
			if tx.To() != nil && tx.To().String() == smartContractAddrFilterString {
				fmt.Println(userAddr, order, tx)
				fmt.Println("Found")
				if tx.Data() != nil {
					calldata := tx.Data()
					if len(calldata) < 4 {
						fmt.Println("invalid calldata")
						continue
					}
					sigdata := calldata[:4]
					method, merr := uswap.Abi.MethodById(sigdata[:4])
					s.Assert().Nil(merr)
					fmt.Println(method.Name)
					argdata := calldata[4:]
					if len(argdata)%32 != 0 {
						fmt.Println("invalid argdata")
						continue
					}
					m := make(map[string]interface{})
					err = method.Inputs.UnpackIntoMap(m, argdata)
					s.Assert().Nil(err)
				}
			}
		}
	}
}

func forceDirToLocation() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	err := os.Chdir(dir)
	if err != nil {
		panic(err.Error())
	}
	return dir
}
func TestWeb3ClientTestSuite(t *testing.T) {
	suite.Run(t, new(Web3ClientTestSuite))
}
