package async_analysis

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
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
	/*
		ALTER TABLE "public"."erc20_token_info"
		ADD COLUMN "name" text,
		ADD COLUMN "symbol" text,
		ADD COLUMN "decimals" int4,
		ADD COLUMN "transfer_tax_numerator" int8,
		ADD COLUMN "transfer_tax_denominator" int8,
		ADD COLUMN "trading_enabled" bool NOT NULL DEFAULT false;
	*/
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
