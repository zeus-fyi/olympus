package artemis_trading_auxiliary

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
)

const (
	ZeusTestSessionLockHeaderValue = "Zeus-Test"
)

func (a *AuxiliaryTradingUtils) setupCleanSimEnvironment(ctx context.Context, bn int) error {
	a.w3c().Web3Actions.AddSessionLockHeader(ZeusTestSessionLockHeaderValue)
	a.w3c().Dial()
	origInfo, err := a.w3c().GetNodeMetadata(ctx)
	if err != nil {
		panic(err)
	}
	a.w3c().Close()
	a.w3c().Dial()
	err = a.w3c().ResetNetwork(ctx, origInfo.ForkConfig.ForkUrl, bn)
	if err != nil {
		panic(err)
	}
	a.w3c().Close()
	eb := artemis_eth_units.EtherMultiple(100000)
	bal := (*hexutil.Big)(eb)
	acc, err := accounts.CreateAccount()
	if err != nil {
		return err
	}
	a.w3c().Account = acc
	a.w3c().Dial()
	defer a.w3c().Close()
	err = a.w3c().SetBalance(ctx, a.w3c().PublicKey(), *bal)
	if err != nil {
		log.Err(err).Msg("error setting balance")
		return err
	}
	approveTx, err := a.w3c().ERC20ApproveSpender(ctx,
		artemis_trading_constants.WETH9ContractAddress,
		artemis_trading_constants.Permit2SmartContractAddress,
		artemis_eth_units.MaxUINT)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}
	err = a.w3c().SetERC20BalanceBruteForce(ctx, artemis_trading_constants.WETH9ContractAddress, a.w3c().PublicKey(), artemis_eth_units.EtherMultiple(10000))
	if err != nil {
		return err
	}
	return err
}
