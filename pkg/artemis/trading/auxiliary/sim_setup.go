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

func (a *AuxiliaryTradingUtils) setupCleanSimEnvironment(ctx context.Context) error {
	a.U.Web3Client.Web3Actions.AddSessionLockHeader(ZeusTestSessionLockHeaderValue)
	a.U.Web3Client.Dial()
	origInfo, err := a.U.Web3Client.GetNodeMetadata(ctx)
	if err != nil {
		panic(err)
	}
	a.U.Web3Client.Close()
	a.U.Web3Client.Dial()
	err = a.U.Web3Client.ResetNetwork(ctx, origInfo.ForkConfig.ForkUrl, 0)
	if err != nil {
		panic(err)
	}
	a.U.Web3Client.Close()
	eb := artemis_eth_units.EtherMultiple(100000)
	bal := (*hexutil.Big)(eb)
	acc, err := accounts.CreateAccount()
	if err != nil {
		return err
	}
	a.U.Web3Client.Account = acc
	a.U.Web3Client.Dial()
	defer a.U.Web3Client.Close()
	err = a.U.Web3Client.SetBalance(ctx, a.U.Web3Client.PublicKey(), *bal)
	if err != nil {
		log.Err(err).Msg("error setting balance")
		return err
	}
	approveTx, err := a.U.Web3Client.ERC20ApproveSpender(ctx,
		artemis_trading_constants.WETH9ContractAddress,
		artemis_trading_constants.Permit2SmartContractAddress,
		artemis_eth_units.MaxUINT)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}
	err = a.U.Web3Client.SetERC20BalanceBruteForce(ctx, artemis_trading_constants.WETH9ContractAddress, a.U.Web3Client.PublicKey(), artemis_eth_units.EtherMultiple(10000))
	if err != nil {
		return err
	}
	return err
}
