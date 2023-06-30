package async_analysis

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_core/entities"

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
	shib2Contract := "0x34ba042827996821CFFEB06477D48a2Ff9474483"
	s.ca = NewERC20ContractAnalysis(&uni, shib2Contract)
	s.ca.UserB = s.UserB
	percent, err := s.ca.CalculateTransferFeeTax(ctx, web3_client.EtherMultiple(1))
	s.Assert().Nil(err)
	s.Assert().Equal(int64(1), percent.Numerator.Int64())
	s.Assert().Equal(int64(50), percent.Denominator.Int64())
	amount := core_entities.Fraction{
		Numerator:   web3_client.Ether,
		Denominator: new(big.Int).SetUint64(1),
	}
	feeAmount := amount.Multiply(percent.Fraction)
	s.Assert().Equal("20000000000000000", feeAmount.Quotient().String())
}

func (s *ArtemisRealTimeTradingTestSuite) TestTransferFeeAnalysisBulk() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	uni := web3_client.InitUniswapClient(ctx, s.UserA)
	uni.Web3Client.IsAnvilNode = true
	shib2Contract := "0x34ba042827996821CFFEB06477D48a2Ff9474483"
	s.ca = NewERC20ContractAnalysis(&uni, shib2Contract)
	s.ca.UserB = s.UserB
	tokens, terr := artemis_validator_service_groups_models.SelectERC20TokensWithoutTransferTaxInfo(ctx)
	s.Assert().Nil(terr)
	s.Assert().NotNil(tokens)

	for _, token := range tokens {
		s.ca.SmartContractAddr = token.Address
		percent, err := s.ca.CalculateTransferFeeTax(ctx, web3_client.EtherMultiple(1))
		s.Assert().Nil(err)
		num := int(percent.Numerator.Int64())
		token.TransferTaxNumerator = &num
		denom := int(percent.Denominator.Int64())
		token.TransferTaxDenominator = &denom
		s.Require().NotZero(token.TransferTaxDenominator)

		fmt.Println("token", token.Address, "percent", percent.Numerator.String(), "/", percent.Denominator.String())
		err = artemis_validator_service_groups_models.UpdateERC20TokenTransferTaxInfo(ctx, token)
		s.Assert().Nil(err)
		time.Sleep(100 * time.Millisecond)
	}
}
func (s *ArtemisRealTimeTradingTestSuite) SetupTest() {
	s.InitLocalConfigs()
	newAccount, err := accounts.ParsePrivateKey("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	s.Assert().Nil(err)

	pkHexString := s.Tc.LocalEcsdaTestPkey
	secondAccount, err := accounts.ParsePrivateKey(pkHexString)
	s.Assert().Nil(err)
	//s.UserA = web3_client.NewWeb3Client(s.Tc.QuiknodeLiveNode, newAccount)
	s.UserA = web3_client.NewWeb3Client("http://localhost:8545", newAccount)
	s.UserB = web3_client.NewWeb3Client("http://localhost:8545", secondAccount)
}

func TestArtemisRealTimeTradingTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisRealTimeTradingTestSuite))
}
