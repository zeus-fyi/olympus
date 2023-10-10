package async_analysis

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"

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
	uni.Web3Client.DurableExecution = true
	shib2Contract := "0x34ba042827996821CFFEB06477D48a2Ff9474483"
	s.ca = NewERC20ContractAnalysis(&uni, shib2Contract)
	s.ca.UserB = s.UserB

	tokens, _, terr := artemis_mev_models.SelectERC20TokensWithNullTransferTax(ctx)
	s.Assert().Nil(terr)
	s.Assert().NotNil(tokens)
	s.ca.UserA.IsAnvilNode = true
	s.ca.UserA.DurableExecution = true
	for _, token := range tokens {
		if token.BalanceOfSlotNum == -1 {
			continue
		}
		s.ca.u.Web3Client.AddSessionLockHeader(token.Address)
		err := s.ca.UserA.HardHatResetNetwork(ctx, 17595510)
		s.Assert().Nil(err)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		fmt.Println("token", token.Address)
		s.ca.SmartContractAddr = token.Address
		percent, err := s.ca.CalculateTransferFeeTax(ctx, artemis_eth_units.EtherMultiple(1))
		s.Assert().Nil(err)
		if err != nil {
			continue
		}

		if percent.Numerator == nil || percent.Denominator == nil {
			num := int(0)
			token.TransferTaxNumerator = &num
			denom := int(1)
			token.TransferTaxDenominator = &denom
			err = artemis_mev_models.UpdateERC20TokenTransferTaxInfo(ctx, token)
			s.Assert().Nil(err)
		} else {
			num := int(percent.Numerator.Int64())
			token.TransferTaxNumerator = &num
			denom := int(percent.Denominator.Int64())
			token.TransferTaxDenominator = &denom
			s.Require().NotZero(token.TransferTaxDenominator)
			fmt.Println("token", token.Address, "percent", percent.Numerator.String(), "/", percent.Denominator.String())
			err = artemis_mev_models.UpdateERC20TokenTransferTaxInfo(ctx, token)
			s.Assert().Nil(err)
		}

		s.Assert().Nil(err)
		time.Sleep(100 * time.Millisecond)
	}
}

func (s *ArtemisRealTimeTradingTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	artemis_trading_cache.InitTokenFilter(ctx)
	//apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	newAccount, err := accounts.ParsePrivateKey("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	s.Assert().Nil(err)

	pkHexString := s.Tc.LocalEcsdaTestPkey
	secondAccount, err := accounts.ParsePrivateKey(pkHexString)
	s.Assert().Nil(err)
	irisBetaSvc := artemis_trading_constants.IrisAnvilRoute

	wc := web3_client.NewWeb3Client(irisBetaSvc, newAccount)
	m := map[string]string{
		"Authorization": "Bearer " + s.Tc.ProductionLocalTemporalBearerToken,
	}
	wc.Headers = m
	wc.AddSessionLockHeader(uuid.New().String())
	uni := web3_client.InitUniswapClient(ctx, wc)
	uni.PrintOn = true
	uni.PrintLocal = false
	uni.Web3Client.IsAnvilNode = true
	uni.Web3Client.DurableExecution = true
	s.UserA = wc
	// web3_client.NewWeb3Client(s.Tc.QuikNodeLiveNode, newAccount)
	//s.UserA = web3_client.NewWeb3Client("http://localhost:8545", newAccount)
	s.UserB = web3_client.NewWeb3Client("http://localhost:8545", secondAccount)
}

func TestArtemisRealTimeTradingTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisRealTimeTradingTestSuite))
}
