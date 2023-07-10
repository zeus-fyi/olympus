package artemis_trading_auxiliary

import "C"
import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (a *AuxiliaryTradingUtils) universalRouterCmdVerifier(ctx context.Context, ur *web3_client.UniversalRouterExecCmd, scInfo *web3_actions.SendContractTxPayload) error {
	ethRequirements := artemis_eth_units.NewBigInt(0)
	for _, sc := range ur.Commands {
		switch sc.Command {
		case artemis_trading_constants.WrapETH:
			wp := sc.DecodedInputs.(web3_client.WrapETHParams)
			ethRequirements = artemis_eth_units.AddBigInt(ethRequirements, wp.AmountMin)
			if ur.Payable == nil {
				return errors.New("payable is nil")
			}
		case artemis_trading_constants.UnwrapWETH:
		}
	}
	gasCost := artemis_eth_units.MulBigInt(scInfo.GasFeeCap, artemis_eth_units.NewBigInt(int(scInfo.GasLimit)))
	fmt.Println("gasCost", gasCost)
	ethRequirements = artemis_eth_units.AddBigInt(ethRequirements, gasCost)
	hasEnough, err := a.checkAuxEthBalanceGreaterThan(ctx, ethRequirements)
	if err != nil {
		return err
	}
	if !hasEnough {
		return errors.New("user does not have enough ETH to exchange to WETH")
	}

	return nil
}

func (a *AuxiliaryTradingUtils) checkAuxEthBalance(ctx context.Context) (*big.Int, error) {
	bal, err := a.u.Web3Client.GetCurrentBalance(ctx)
	if err != nil {
		return bal, err
	}
	return bal, err
}

func (a *AuxiliaryTradingUtils) getAccountAddressString() string {
	return a.getWeb3Client().Account.Address().String()
}

func (a *AuxiliaryTradingUtils) checkAuxERC20Balance(ctx context.Context, token core_entities.Token) (*big.Int, error) {
	bal, err := a.getWeb3Client().ReadERC20TokenBalance(ctx, token.Address.String(), a.getAccountAddressString())
	if err != nil {
		return bal, err
	}
	return bal, err
}

func (a *AuxiliaryTradingUtils) getWeb3Client() *web3_client.Web3Client {
	return &a.u.Web3Client
}

func (a *AuxiliaryTradingUtils) checkAuxERC20BalanceGreaterThan(ctx context.Context, token core_entities.Token, amount *big.Int) (bool, error) {
	bal, err := a.checkAuxERC20Balance(ctx, token)
	if err != nil {
		return false, err
	}
	return artemis_eth_units.IsXGreaterThanY(bal, amount), err
}

func (a *AuxiliaryTradingUtils) checkAuxEthBalanceGreaterThan(ctx context.Context, amount *big.Int) (bool, error) {
	bal, err := a.getWeb3Client().GetCurrentBalance(ctx)
	if err != nil {
		return false, err
	}
	return artemis_eth_units.IsXGreaterThanY(bal, amount), err
}

func (a *AuxiliaryTradingUtils) checkAuxWETHBalance(ctx context.Context) (*big.Int, error) {
	wethAddr := a.getChainSpecificWETH()
	chainID, err := a.getChainID(ctx)
	if err != nil {
		return nil, err
	}
	token := core_entities.NewToken(uint(chainID), wethAddr, 18, "WETH", "Wrapped Ether")

	bal, err := a.getWeb3Client().ReadERC20TokenBalance(ctx, token.Address.String(), a.getAccountAddressString())
	if err != nil {
		return bal, err
	}
	return bal, err
}

func (a *AuxiliaryTradingUtils) checkAuxWETHBalanceGreaterThan(ctx context.Context, amount *big.Int) (bool, error) {
	bal, err := a.checkAuxWETHBalance(ctx)
	if err != nil {
		return false, err
	}
	return artemis_eth_units.IsXGreaterThanY(bal, amount), err
}

func (a *AuxiliaryTradingUtils) getChainID(ctx context.Context) (int, error) {
	chainID := hestia_req_types.EthereumMainnetProtocolNetworkID
	switch a.getWeb3Client().Network {
	case hestia_req_types.Mainnet:
		chainID = hestia_req_types.EthereumMainnetProtocolNetworkID
	case hestia_req_types.Goerli:
		chainID = hestia_req_types.EthereumGoerliProtocolNetworkID
	case hestia_req_types.Ephemery:
		chainID = hestia_req_types.EthereumEphemeryProtocolNetworkID
	default:
		a.getWeb3Client().Dial()
		chain, err := a.getWeb3Client().C.ChainID(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("error getting chainID")
			return 0, err
		}
		chainID = int(chain.Int64())
		a.getWeb3Client().Close()
	}
	return chainID, nil
}

func (a *AuxiliaryTradingUtils) getChainSpecificWETH() accounts.Address {
	wethAddr := artemis_trading_constants.WETH9ContractAddressAccount
	switch a.getWeb3Client().Network {
	case hestia_req_types.Mainnet:
		wethAddr = artemis_trading_constants.WETH9ContractAddressAccount
	case hestia_req_types.Goerli:
		wethAddr = artemis_trading_constants.GoerliWETH9ContractAddressAccount
	}
	return wethAddr
}
