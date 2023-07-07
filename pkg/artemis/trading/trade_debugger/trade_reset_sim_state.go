package artemis_trade_debugger

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog/log"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (t *TradeDebugger) ResetAndSetupPreconditions(ctx context.Context, tf web3_client.TradeExecutionFlow) error {
	err := t.resetNetwork(ctx, tf)
	if err != nil {
		return err
	}
	err = t.setupCleanEnvironment(ctx, tf)
	if err != nil {
		return err
	}
	err = t.UniswapClient.CheckExpectedReserves(&tf)
	if err != nil {
		return err
	}
	err = t.UniswapClient.Web3Client.MatchFrontRunTradeValues(&tf)
	if err != nil {
		return err
	}
	return nil
}

func (t *TradeDebugger) resetNetwork(ctx context.Context, tf web3_client.TradeExecutionFlow) error {
	if tf.CurrentBlockNumber == nil {
		return fmt.Errorf("current block number is nil")
	}
	bn, err := t.UniswapClient.CheckBlockRxAndNetworkReset(ctx, &tf, &t.LiveNetworkClient)
	if err != nil {
		log.Err(err).Interface("blockNum", bn).Msg("error checking block and network reset")
		return err
	}

	// FOR DEBUGGING: uncomment this to check if the sim block num is the same as the live block num
	simBlockNum, err := t.UniswapClient.Web3Client.GetBlockHeight(ctx)
	if err != nil {
		return err
	}
	nodeInfo, err := t.UniswapClient.Web3Client.GetNodeMetadata(ctx)
	if err != nil {
		return err
	}

	if nodeInfo.ForkConfig.ForkBlockNumber != bn {
		return fmt.Errorf("sim block num %s != live block num %d", simBlockNum, bn)
	}
	return err
}

func (t *TradeDebugger) setupCleanEnvironment(ctx context.Context, tf web3_client.TradeExecutionFlow) error {
	eb := artemis_eth_units.EtherMultiple(10000)
	bal := (*hexutil.Big)(eb)
	t.UniswapClient.Web3Client.Dial()
	defer t.UniswapClient.Web3Client.Close()
	err := t.UniswapClient.Web3Client.SetBalance(ctx, t.UniswapClient.Web3Client.PublicKey(), *bal)
	if err != nil {
		log.Err(err).Msg("error setting balance")
		return err
	}
	nv, _ := new(big.Int).SetString("0", 10)
	nvB := (*hexutil.Big)(nv)
	err = t.UniswapClient.Web3Client.SetNonce(ctx, t.UniswapClient.Web3Client.PublicKey(), *nvB)
	if err != nil {
		log.Err(err).Msg("error setting nonce")
		return err
	}
	approveTx, err := t.UniswapClient.ApproveSpender(ctx,
		artemis_trading_constants.WETH9ContractAddress,
		artemis_trading_constants.Permit2SmartContractAddress,
		artemis_eth_units.MaxUINT)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}
	secondToken := tf.FrontRunTrade.AmountInAddr.String()
	if tf.FrontRunTrade.AmountInAddr.String() == artemis_trading_constants.WETH9ContractAddress {
		secondToken = tf.FrontRunTrade.AmountOutAddr.String()
	}
	approveTx, err = t.UniswapClient.ApproveSpender(ctx, secondToken, artemis_trading_constants.Permit2SmartContractAddress, artemis_eth_units.MaxUINT)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}
	return err
}
