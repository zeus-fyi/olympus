package web3_client

import (
	"encoding/json"
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (s *Web3ClientTestSuite) TestTradeExec() {
	//forceDirToLocation()
	//swapAbi, bc, err := LoadSwapAbiPayload()
	//s.Require().NoError(err)
	//s.Require().NotNil(swapAbi)
	//s.Require().NotNil(bc)
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	//uni := InitUniswapV2Client(ctx, s.LocalHardhatMainnetUser)
	mevTxs, merr := artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17326677)
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)

	for _, mevTx := range mevTxs {
		tf := TradeExecutionFlow{}
		by := []byte(mevTx.TxFlowPrediction)
		berr := json.Unmarshal(by, &tf)
		s.Require().Nil(berr)

		if tf.FrontRunTrade.AmountIn == "" {
			continue
		}

		tfRegular := tf.ConvertToBigIntType()
		addrIn := tfRegular.FrontRunTrade.AmountInAddr.String()
		b, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, addrIn, s.LocalHardhatMainnetUser.PublicKey())
		s.Require().Nil(err)
		s.Assert().NotZero(b)
		fmt.Println(b.String())

		err = s.LocalHardhatMainnetUser.SetERC20BalanceBruteForce(ctx, tf.FrontRunTrade.AmountInAddr.String(), s.LocalHardhatMainnetUser.PublicKey(), tfRegular.FrontRunTrade.AmountIn)
		s.Require().Nil(err)
		b, err = s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, tf.FrontRunTrade.AmountInAddr.String(), s.LocalHardhatMainnetUser.PublicKey())
		s.Require().Nil(err)
		s.Assert().NotZero(b)
		fmt.Println(b.String())
		s.Assert().Equal(tfRegular.FrontRunTrade.AmountIn.String(), b.String())
	}
}

func (s *Web3ClientTestSuite) TestMatchInputs() {
	rerr := s.LocalHardhatMainnetUser.ResetNetwork(ctx, s.Tc.HardhatNode, 17326550)
	s.Require().Nil(rerr)
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	//uni := InitUniswapV2Client(ctx, s.LocalHardhatMainnetUser)
	mevTxs, merr := artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17326677)
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)

	for _, mevTx := range mevTxs {
		tf := TradeExecutionFlow{}
		by := []byte(mevTx.TxFlowPrediction)
		berr := json.Unmarshal(by, &tf)
		s.Require().Nil(berr)

		if tf.FrontRunTrade.AmountIn == "" {
			continue
		}

		tfRegular := tf.ConvertToBigIntType()
		err := s.LocalHardhatMainnetUser.MatchFrontRunTradeValues(tfRegular)
		s.Require().Nil(err)
		b, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, tf.FrontRunTrade.AmountInAddr.String(), s.LocalHardhatMainnetUser.PublicKey())
		s.Require().Nil(err)
		s.Require().Equal(tfRegular.FrontRunTrade.AmountIn.String(), b.String())
		fmt.Println(b.String(), tfRegular.FrontRunTrade.AmountIn.String())

		uni := InitUniswapV2Client(ctx, s.LocalMainnetWeb3User)
		aa, err := uni.SwapExactTokensForETHParams(tfRegular)
		s.Require().Nil(err)
		s.Require().NotNil(aa)
	}
}
