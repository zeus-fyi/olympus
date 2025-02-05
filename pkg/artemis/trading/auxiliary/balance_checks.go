package artemis_trading_auxiliary

import "C"
import (
	"context"
	"errors"
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

func universalRouterCmdVerifier(ctx context.Context, w3c web3_client.Web3Client, ur *web3_client.UniversalRouterExecCmd, scInfo *web3_actions.SendContractTxPayload) error {
	ethRequirements := artemis_eth_units.NewBigInt(0)
	for _, sc := range ur.Commands {
		switch sc.Command {
		case artemis_trading_constants.WrapETH:
			wp := sc.DecodedInputs.(web3_client.WrapETHParams)
			ethRequirements = artemis_eth_units.AddBigInt(ethRequirements, wp.AmountMin)
			if ur.Payable == nil {
				log.Warn().Msgf("universalRouterCmdVerifier: payable is nil")
				return errors.New("payable is nil")
			}
		case artemis_trading_constants.UnwrapWETH:
		}
	}
	gasCost := artemis_eth_units.MulBigInt(scInfo.GasFeeCap, artemis_eth_units.NewBigInt(int(scInfo.GasLimit)))
	ethRequirements = artemis_eth_units.AddBigInt(ethRequirements, gasCost)
	hasEnough, err := checkEthBalanceGreaterThan(ctx, w3c, ethRequirements)
	if err != nil {
		log.Err(err).Str("ethRequirements", ethRequirements.String()).Msgf("universalRouterCmdVerifier: checkEthBalanceGreaterThan checking ETH balance")
		return err
	}
	if !hasEnough {
		log.Warn().Msgf("universalRouterCmdVerifier: user does not have enough ETH to execute the command")
		return errors.New("user does not have enough ETH to exchange to WETH")
	}

	return nil
}

func checkEthBalance(ctx context.Context, w3c web3_client.Web3Client) (*big.Int, error) {
	bal, err := w3c.GetCurrentBalance(ctx)
	if err != nil {
		return bal, err
	}
	return bal, err
}

func (a *AuxiliaryTradingUtils) getAccountAddressString() string {
	return a.w3c().Account.Address().String()
}

func (a *AuxiliaryTradingUtils) checkAuxERC20Balance(ctx context.Context, token core_entities.Token) (*big.Int, error) {
	bal, err := a.w3c().ReadERC20TokenBalance(ctx, token.Address.String(), a.getAccountAddressString())
	if err != nil {
		return bal, err
	}
	return bal, err
}

func (a *AuxiliaryTradingUtils) CheckAuxERC20BalanceFromAddr(ctx context.Context, token string) (*big.Int, error) {
	bal, err := a.w3c().ReadERC20TokenBalance(ctx, token, a.getAccountAddressString())
	if err != nil {
		return bal, err
	}
	return bal, err
}

func (a *AuxiliaryTradingUtils) checkAuxERC20BalanceGreaterThan(ctx context.Context, token core_entities.Token, amount *big.Int) (bool, error) {
	bal, err := a.checkAuxERC20Balance(ctx, token)
	if err != nil {
		return false, err
	}
	return artemis_eth_units.IsXGreaterThanY(bal, amount), err
}

func checkEthBalanceGreaterThan(ctx context.Context, w3c web3_client.Web3Client, amount *big.Int) (bool, error) {
	bal, err := w3c.GetCurrentBalance(ctx)
	if err != nil {
		return false, err
	}
	log.Info().Msgf("ETH balance: %s", bal.String())
	return artemis_eth_units.IsXGreaterThanY(bal, amount), err
}

func CheckEthBalanceGreaterThan(ctx context.Context, w3c web3_client.Web3Client, amount *big.Int) (bool, error) {
	bal, err := w3c.GetCurrentBalance(ctx)
	if err != nil {
		return false, err
	}
	log.Info().Msgf("ETH balance: %s", bal.String())
	return artemis_eth_units.IsXGreaterThanY(bal, amount), err
}

func CheckAuxWETHBalance(ctx context.Context, w3c web3_client.Web3Client) (*big.Int, error) {
	wethAddr := getChainSpecificWETH(w3c)
	chainID, err := getChainID(ctx, w3c)
	if err != nil {
		return nil, err
	}
	token := core_entities.NewToken(uint(chainID), wethAddr, 18, "WETH", "Wrapped Ether")
	bal, err := w3c.ReadERC20TokenBalance(ctx, token.Address.String(), w3c.Account.PublicKey())
	if err != nil {
		return bal, err
	}
	return bal, err
}

func CheckAuxWETHBalanceGreaterThan(ctx context.Context, w3c web3_client.Web3Client, amount *big.Int) (bool, error) {
	bal, err := CheckAuxWETHBalance(ctx, w3c)
	if err != nil {
		return false, err
	}
	return artemis_eth_units.IsXGreaterThanY(bal, amount), err
}

func CheckMainnetAuxWETHBalanceGreaterThan(ctx context.Context, w3c web3_client.Web3Client, amount *big.Int) (bool, error) {
	token := core_entities.NewToken(uint(1), artemis_trading_constants.WETH9ContractAddressAccount, 18, "WETH", "Wrapped Ether")
	bal, err := w3c.ReadERC20TokenBalance(ctx, token.Address.String(), w3c.Account.Address().String())
	if err != nil {
		log.Warn().Err(err).Msg("error getting WETH balance")
		return false, err
	}
	return artemis_eth_units.IsXGreaterThanY(bal, amount), err
}
func getChainID(ctx context.Context, w3c web3_client.Web3Client) (int, error) {
	chainID := hestia_req_types.EthereumMainnetProtocolNetworkID
	switch w3c.Network {
	case hestia_req_types.Mainnet:
		chainID = hestia_req_types.EthereumMainnetProtocolNetworkID
	case hestia_req_types.Goerli:
		chainID = hestia_req_types.EthereumGoerliProtocolNetworkID
	case hestia_req_types.Ephemery:
		chainID = hestia_req_types.EthereumEphemeryProtocolNetworkID
	default:
		w3c.Dial()
		defer w3c.Close()
		chain, err := w3c.C.ChainID(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("error getting chainID")
			return -1, err
		}
		chainID = int(chain.Int64())
		return chainID, nil
	}
	return chainID, nil
}

func getChainSpecificWETH(w3c web3_client.Web3Client) accounts.Address {
	wethAddr := artemis_trading_constants.WETH9ContractAddressAccount
	switch w3c.Network {
	case hestia_req_types.Mainnet:
		wethAddr = artemis_trading_constants.WETH9ContractAddressAccount
	case hestia_req_types.Goerli:
		wethAddr = artemis_trading_constants.GoerliWETH9ContractAddressAccount
	default:
		wethAddr = artemis_trading_constants.WETH9ContractAddressAccount
	}
	return wethAddr
}
