package web3_client

import (
	"encoding/json"
	"fmt"
	"math/big"

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
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	mevTxs, merr := artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17332397)
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
		err := s.LocalHardhatMainnetUser.HardhatResetNetworkToBlockBeforeTxMined(ctx, s.Tc.HardhatNode, s.LocalHardhatMainnetUser, s.MainnetWeb3User, *tf.Tx.Hash)
		s.Require().Nil(err)

		tfRegular := tf.ConvertToBigIntType()
		err = s.LocalHardhatMainnetUser.MatchFrontRunTradeValues(tfRegular)
		s.Require().Nil(err)

		uni := InitUniswapV2Client(ctx, s.LocalHardhatMainnetUser)
		uni.DebugPrint = true
		_, _ = uni.FrontRunTradeGetAmountsOut(tfRegular)
		_, _ = uni.UserTradeGetAmountsOut(tfRegular)
		_, err = uni.ExecFrontRunTradeStepTokenTransfer(&tfRegular)
		s.Require().Nil(err)

		_, _ = uni.UserTradeGetAmountsOut(tfRegular)
		// must exceed 33073549076721602

		startBal, err := s.LocalHardhatMainnetUser.GetBalance(ctx, tfRegular.Tx.From.String(), nil)
		fmt.Println("userTradeStartEthBal", startBal.String())

		aa, err := uni.ExecUserTradeStep(&tfRegular)
		s.Require().Nil(err)
		s.Require().NotNil(aa)

		fmt.Println(tfRegular.Tx.Hash.String())
		ethBalance, err := s.LocalHardhatMainnetUser.GetBalance(ctx, tfRegular.Tx.From.String(), nil)
		s.Require().Nil(err)
		s.Require().NotNil(ethBalance)
		fmt.Println("userTradeAmountOut", ethBalance.String())
		s.Require().Nil(err)

		gasUsed := new(big.Int).SetInt64(36627061988 * 114409)
		balanceDiff := new(big.Int).Sub(ethBalance, startBal)
		balanceDiff = new(big.Int).Add(balanceDiff, gasUsed)
		fmt.Println("balanceDiff", balanceDiff.String())

		_, _ = uni.SandwichTradeGetAmountsOut(tfRegular)
		fmt.Println("sandwich expected amounts in/out", tfRegular.SandwichTrade.AmountIn.String(), tfRegular.SandwichTrade.AmountOut.String())
		s.Require().Nil(err)

		_, err = uni.ExecSandwichTradeStepTokenTransfer(&tfRegular)
		s.Require().Nil(err)

		err = tfRegular.GetAggregateGasUsage(ctx, s.LocalHardhatMainnetUser)
		s.Require().Nil(err)
	}
}
