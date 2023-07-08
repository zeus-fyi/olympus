package async_analysis

import (
	"math/big"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (s *ArtemisRealTimeTradingTestSuite) TestEthSimTransferFeeAnalysis() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	uni := web3_client.InitUniswapClient(ctx, s.UserA)
	shib2Contract := "0x34ba042827996821CFFEB06477D48a2Ff9474483"
	s.ca = NewERC20ContractAnalysis(&uni, shib2Contract)
	s.ca.UserB = s.UserB
	err := s.ca.UserA.HardHatResetNetwork(ctx, 17601900)
	s.Require().Nil(err)
	transferTaxPercent, err := s.ca.CalculateTransferFeeTax(ctx, web3_client.EtherMultiple(1))
	s.Assert().Nil(err)
	s.Assert().Equal(int64(1), transferTaxPercent.Numerator.Int64())
	s.Assert().Equal(int64(50), transferTaxPercent.Denominator.Int64())

	// this isn't included in trade gas costs since we amortize one time gas costs for permit2
	max, _ := new(big.Int).SetString(web3_client.MaxUINT, 10)
	approveTx, err := s.ca.u.ApproveSpender(ctx, artemis_trading_constants.WETH9ContractAddressAccount.String(), web3_client.Permit2SmartContractAddress, max)
	s.Assert().Nil(err)
	s.Assert().NotNil(approveTx)

	approveTx, err = s.ca.u.ApproveSpender(ctx, artemis_trading_constants.WETH9ContractAddressAccount.String(), "0x34ba042827996821CFFEB06477D48a2Ff9474483", max)
	s.Assert().Nil(err)
	s.Assert().NotNil(approveTx)
	_, err = s.ca.SimEthTransferFeeTaxTrade(ctx, web3_client.EtherMultiple(1))
	s.Assert().Nil(err)
}
