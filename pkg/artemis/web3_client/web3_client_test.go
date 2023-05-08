package web3_client

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/gochain/gochain/v4/common"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/signing_automation/ethereum"
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

	s.MainnetWeb3User = NewWeb3Client(s.Tc.MainnetNodeUrl, newAccount)
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
	forceDirToLocation()
	s.MainnetWeb3User.Web3Actions.Dial()
	defer s.MainnetWeb3User.Close()
	mempool, err := s.MainnetWeb3User.Web3Actions.GetTxPoolContent(ctx)
	s.Require().Nil(err)
	s.Assert().NotNil(mempool)
	//abiConst := artemis_oly_contract_abis.UniswapV2RouterABI
	addrFilter := "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	abiFile, err := signing_automation_ethereum.ABIOpenFile(ctx, "./contract_abis/uniswapV2Abi.json")
	s.Require().Nil(err)

	smartContractAddrFilter := common.HexToAddress(addrFilter)
	smartContractAddrFilterString := smartContractAddrFilter.String()
	for userAddr, txPoolQueue := range mempool["pending"] {
		for order, tx := range txPoolQueue {
			if tx.To != nil && tx.To.String() == smartContractAddrFilterString {
				fmt.Println(userAddr, order, tx)
				fmt.Println("Found")
				if tx.Input != nil {
					input := *tx.Input
					calldata := []byte(input)
					sigdata := calldata[:4]
					method, merr := abiFile.MethodById(sigdata[:4])
					s.Assert().Nil(merr)
					fmt.Println(method.Name)
					argdata := calldata[4:]
					if len(argdata)%32 != 0 {
						fmt.Println("invalid argdata")
						continue
					}
					inputValues, uerr := method.Inputs.UnpackValues(argdata)
					s.Assert().Nil(uerr)
					fmt.Println(inputValues)
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
