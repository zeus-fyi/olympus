package artemis_realtime_trading

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) setupCleanSimEnvironment(ctx context.Context, tf *web3_client.TradeExecutionFlow) error {
	eb := artemis_eth_units.EtherMultiple(100000)
	bal := (*hexutil.Big)(eb)
	acc, err := accounts.CreateAccount()
	if err != nil {
		return err
	}
	a.simW3c().Account = acc
	a.simW3c().Dial()
	defer a.simW3c().Close()
	err = a.simW3c().SetBalance(ctx, a.simW3c().PublicKey(), *bal)
	if err != nil {
		log.Err(err).Msg("error setting balance")
		return err
	}
	nv, _ := new(big.Int).SetString("0", 10)
	nvB := (*hexutil.Big)(nv)
	err = a.simW3c().SetNonce(ctx, a.simW3c().PublicKey(), *nvB)
	if err != nil {
		log.Err(err).Msg("error setting nonce")
		return err
	}
	approveTx, err := a.simW3c().ERC20ApproveSpender(ctx,
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
	approveTx, err = a.simW3c().ERC20ApproveSpender(ctx, secondToken, artemis_trading_constants.Permit2SmartContractAddress, artemis_eth_units.MaxUINT)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}
	return err
}
