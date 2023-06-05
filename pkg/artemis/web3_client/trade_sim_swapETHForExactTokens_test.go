package web3_client

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

/*
	blockNum 17383905
	tradeMethod swapETHForExactTokens
	txHash 0x7666aee2aeef8c6a069c3e4204d4ccf7462a15641df01f5f806ce6f40860d947
	rxBlockNum 17383906
	blockNum recorded from artemis 17383905

ERRORS

blockNum 17384004
tradeMethod swapETHForExactTokens
txHash 0x8fd935462b382f20133824263d5a598bc71087714b515ae5400b68d062acdc30
{"level":"error","error":"Error: VM Exception while processing transaction: reverted with reason string 'Insufficient Balance'","time":"2023-05-31T22:53:04-07:00","message":"error executing sandwich trade step token transfer"}

*/

func (s *Web3ClientTestSuite) TestFullSandwichTradeSim_SwapETHForExactTokens() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	mevTxs, merr := artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17383905)
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)
	for _, mevTx := range mevTxs {
		tf := TradeExecutionFlow{}
		by := []byte(mevTx.TxFlowPrediction)
		berr := json.Unmarshal(by, &tf)
		s.Require().Nil(berr)
		if tf.FrontRunTrade.AmountIn == "0" || tf.FrontRunTrade.AmountIn == "" {
			continue
		}
		if tf.Trade.TradeMethod != swapETHForExactTokens {
			continue
		}
		s.Require().Equal(swapETHForExactTokens, tf.Trade.TradeMethod)
		fmt.Println("blockNum recorded from artemis", tf.CurrentBlockNumber)
		currentBlockStr := tf.CurrentBlockNumber.String()
		currentBlockNum, err := strconv.Atoi(currentBlockStr)
		s.Require().Nil(err)
		fmt.Println("blockNum recorded from artemis", currentBlockNum)
		err = s.HostedHardhatMainnetUser.HardHatResetNetwork(ctx, s.Tc.HardhatNode, currentBlockNum)
		s.Require().Nil(err)

		tfRegular := tf.ConvertToBigIntType()
		uni := InitUniswapClient(ctx, s.HostedHardhatMainnetUser)
		pairAddr := tfRegular.InitialPair.PairContractAddr
		simPair, err := uni.GetPairContractPrices(ctx, pairAddr)
		s.Require().Nil(err)
		s.Require().Equal(tfRegular.InitialPair.Reserve0.String(), simPair.Reserve0.String())
		s.Require().Equal(tfRegular.InitialPair.Reserve1.String(), simPair.Reserve1.String())

		uni.DebugPrint = true
		err = uni.SimFullSandwichTrade(&tfRegular)
		s.Require().Nil(err)
	}
}
