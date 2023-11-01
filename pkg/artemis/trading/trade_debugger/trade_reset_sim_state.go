package artemis_trade_debugger

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

const (
	ZeusTestSessionLockHeaderValue = "Zeus-Test"
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
	err = t.dat.GetSimUniswapClient().CheckExpectedReserves(&tf)
	if err != nil {
		return err
	}
	err = t.dat.SimW3c().MatchFrontRunTradeValues(&tf)
	if err != nil {
		return err
	}
	return nil
}

func (t *TradeDebugger) resetNetwork(ctx context.Context, tf web3_client.TradeExecutionFlow) error {
	if tf.CurrentBlockNumber == nil {
		return fmt.Errorf("current block number is nil")
	}
	sessionID := t.dat.GetSimUniswapClient().Web3Client.GetSessionLockHeader()
	if t.dat.GetSimUniswapClient().Web3Client.GetSessionLockHeader() == "" {
		t.dat.GetSimUniswapClient().Web3Client.AddSessionLockHeader(ZeusTestSessionLockHeaderValue)
	}
	bn, err := t.dat.GetSimUniswapClient().CheckBlockRxAndNetworkReset(ctx, &tf)
	if err != nil {
		log.Err(err).Interface("blockNum", bn).Str("sessionID", sessionID).Msg("error checking block and network reset")
		return err
	}

	// FOR DEBUGGING: uncomment this to check if the sim block num is the same as the live block num
	//simBlockNum, err := t.dat.SimW3c().GetBlockHeight(ctx)
	//if err != nil {
	//	return err
	//}
	//nodeInfo, err := t.dat.SimW3c().GetNodeMetadata(ctx)
	//if err != nil {
	//	return err
	//}
	//
	//if nodeInfo.ForkConfig.ForkBlockNumber != bn {
	//	return fmt.Errorf("sim block num %s != live block num %d", simBlockNum, bn)
	//}
	return err
}

func (t *TradeDebugger) setupCleanEnvironment(ctx context.Context, tf web3_client.TradeExecutionFlow) error {
	eb := artemis_eth_units.EtherMultiple(100000)
	bal := (*hexutil.Big)(eb)
	acc, err := accounts.CreateAccount()
	if err != nil {
		return err
	}
	if tf.Tx == nil {
		return fmt.Errorf("tx is nil")
	}
	if t.dat.GetSimUniswapClient().Web3Client.GetSessionLockHeader() == "" {
		t.dat.SimW3c().AddSessionLockHeader(ZeusTestSessionLockHeaderValue)
	}
	t.dat.SimW3c().Account = acc
	t.dat.SimW3c().Dial()
	defer t.dat.SimW3c().Close()
	err = t.dat.SimW3c().SetBalance(ctx, t.dat.SimW3c().PublicKey(), *bal)
	if err != nil {
		log.Err(err).Msg("error setting balance")
		return err
	}
	nv, _ := new(big.Int).SetString("0", 10)
	nvB := (*hexutil.Big)(nv)
	err = t.dat.SimW3c().SetNonce(ctx, t.dat.SimW3c().PublicKey(), *nvB)
	if err != nil {
		log.Err(err).Msg("error setting nonce")
		return err
	}
	approveTx, err := t.dat.GetSimUniswapClient().ApproveSpender(ctx,
		artemis_trading_constants.WETH9ContractAddress,
		artemis_trading_constants.Permit2SmartContractAddress,
		artemis_eth_units.MaxUINT)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}
	approveTx, err = t.dat.GetSimUniswapClient().ApproveSpender(ctx, tf.UserTrade.AmountInAddr.String(), artemis_trading_constants.Permit2SmartContractAddress, artemis_eth_units.MaxUINT)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}
	approveTx, err = t.dat.GetSimUniswapClient().ApproveSpender(ctx, tf.UserTrade.AmountOutAddr.String(), artemis_trading_constants.Permit2SmartContractAddress, artemis_eth_units.MaxUINT)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}
	return err
}
