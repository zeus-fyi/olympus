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
		txHash := *tf.Tx.Hash
		fmt.Println(txHash.String())
		s.MainnetWeb3User.Dial()
		rx, err := s.MainnetWeb3User.GetTransactionByHash(ctx, txHash)
		s.MainnetWeb3User.Close()
		s.Require().Nil(err)
		err = s.LocalHardhatMainnetUser.ResetNetwork(ctx, s.Tc.HardhatNode, int(rx.BlockNumber.Int64()-1))
		s.Require().Nil(err)

		ethBalance, err := s.LocalHardhatMainnetUser.GetBalance(ctx, s.LocalHardhatMainnetUser.PublicKey(), nil)
		s.Require().Nil(err)
		s.Require().NotNil(ethBalance)
		fmt.Println("tradeUserEthBalance", ethBalance.String())
		tfRegular := tf.ConvertToBigIntType()
		err = s.LocalHardhatMainnetUser.MatchFrontRunTradeValues(tfRegular)
		s.Require().Nil(err)
		b, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, tf.FrontRunTrade.AmountInAddr.String(), s.LocalHardhatMainnetUser.PublicKey())
		s.Require().Nil(err)
		s.Require().Equal(tfRegular.FrontRunTrade.AmountIn.String(), b.String())
		fmt.Println("frontRunAmountIn", b.String(), tfRegular.FrontRunTrade.AmountIn.String())
		//
		uni := InitUniswapV2Client(ctx, s.LocalHardhatMainnetUser)
		//amounts, err := uni.FrontRunTradeGetAmountsOut(tfRegular)
		//s.Require().Nil(err)
		//s.Require().NotEmpty(amounts)
		//s.Require().Len(amounts, 2)
		//s.Assert().Equal(tfRegular.FrontRunTrade.AmountIn.String(), amounts[0].String())
		//s.Assert().Equal(tfRegular.FrontRunTrade.AmountOut.String(), amounts[1].String())
		//
		//fmt.Println("amountIn", tfRegular.FrontRunTrade.AmountInAddr.String(), tfRegular.FrontRunTrade.AmountIn.String())
		//fmt.Println("amountOut", tfRegular.FrontRunTrade.AmountOutAddr.String(), tfRegular.FrontRunTrade.AmountOut.String())
		//
		amounts, err := uni.UserTradeGetAmountsOut(tfRegular)
		fmt.Println("user expected amounts without front run", amounts[0].String(), amounts[1].String())

		err = uni.RouterApproveAndSend(ctx, tfRegular.FrontRunTrade, tfRegular.InitialPair.PairContractAddr)
		s.Require().Nil(err)
		userTradeMethod := tfRegular.Trade.TradeMethod
		tfRegular.Trade.TradeMethod = swapFrontRun

		out, err := uni.ExecTradeByMethod(tfRegular)
		s.Require().Nil(err)
		s.Require().NotNil(out)

		b, err = s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, tf.FrontRunTrade.AmountOutAddr.String(), s.LocalHardhatMainnetUser.PublicKey())
		s.Require().Nil(err)
		s.Require().Equal(tfRegular.FrontRunTrade.AmountOut.String(), b.String())

		amounts, err = uni.UserTradeGetAmountsOut(tfRegular)
		fmt.Println("user expected amounts post front run", amounts[0].String(), amounts[1].String())
		// must exceed 33073549076721602

		startBal, err := s.LocalHardhatMainnetUser.GetBalance(ctx, tfRegular.Tx.From.String(), nil)
		fmt.Println("userTradeStartEthBal", startBal.String())

		tfRegular.Trade.TradeMethod = userTradeMethod
		aa, err := uni.ExecTradeByMethod(tfRegular)
		s.Require().Nil(err)
		s.Require().NotNil(aa)

		fmt.Println("userAddr", tfRegular.Tx.From.String())
		fmt.Println("amountOutAddr", tf.UserTrade.AmountOutAddr.String())

		ethBalance, err = s.LocalHardhatMainnetUser.GetBalance(ctx, tfRegular.Tx.From.String(), nil)
		s.Require().Nil(err)
		s.Require().NotNil(ethBalance)
		fmt.Println("userTradeAmountOut", ethBalance.String())
		s.Require().Nil(err)

		gasUsed := new(big.Int).SetInt64(36627061988 * 114409)
		balanceDiff := new(big.Int).Sub(ethBalance, startBal)
		balanceDiff = new(big.Int).Add(balanceDiff, gasUsed)
		fmt.Println("balanceDiff", balanceDiff.String())

	}
}
