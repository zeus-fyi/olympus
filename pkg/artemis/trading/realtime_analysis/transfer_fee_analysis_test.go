package artemis_realtime_analysis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"

	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
)

var ctx = context.Background()

type ArtemisRealTimeTradingTestSuite struct {
	test_suites_encryption.EncryptionTestSuite
	ca    ContractAnalysis
	UserA web3_client.Web3Client
	UserB web3_client.Web3Client
}

func (s *ArtemisRealTimeTradingTestSuite) TestTransferFeeAnalysis() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	uni := web3_client.InitUniswapClient(ctx, s.UserA)

	abiFile := web3_client.MustLoadERC20Abi()
	shib2Contract := "0x34ba042827996821CFFEB06477D48a2Ff9474483"
	s.ca = NewContractAnalysis(&uni, shib2Contract, abiFile)
	s.ca.UserB = s.UserB

	percent, err := s.ca.CalculateTransferFeeTax(ctx, web3_client.EtherMultiple(1))
	s.Assert().Nil(err)
	s.Assert().Equal(int64(2), percent)
}

func (s *ArtemisRealTimeTradingTestSuite) SetupTest() {
	s.InitLocalConfigs()
	newAccount, err := accounts.ParsePrivateKey("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	s.Assert().Nil(err)

	pkHexString := s.Tc.LocalEcsdaTestPkey
	secondAccount, err := accounts.ParsePrivateKey(pkHexString)
	s.Assert().Nil(err)
	s.UserA = web3_client.NewWeb3Client("http://localhost:8545", newAccount)
	s.UserB = web3_client.NewWeb3Client("http://localhost:8545", secondAccount)

}

func TestArtemisRealTimeTradingTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisRealTimeTradingTestSuite))
}
