package async_analysis

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	uniswap_core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
)

func (c *ContractAnalysis) SimEthTransferFeeTaxTrade(ctx context.Context, amount *big.Int) (*uniswap_core_entities.Percent, error) {
	err := c.UserA.SetERC20BalanceBruteForce(ctx, artemis_trading_constants.WETH9ContractAddressAccount.String(), c.UserA.Address().String(), amount)
	if err != nil {
		return nil, err
	}
	pd, err := c.u.GetV2PricingData(ctx, []accounts.Address{artemis_trading_constants.WETH9ContractAddressAccount, accounts.HexToAddress(c.SmartContractAddr)})
	if err != nil {
		return nil, err
	}
	calculatedOut, err := pd.V2Pair.PriceImpact(artemis_trading_constants.WETH9ContractAddressAccount, amount)
	if err != nil {
		return nil, err
	}
	trade := &artemis_trading_types.TradeOutcome{
		AmountIn:      amount,
		AmountInAddr:  artemis_trading_constants.WETH9ContractAddressAccount,
		AmountOut:     calculatedOut.AmountOut,
		AmountOutAddr: accounts.HexToAddress(c.SmartContractAddr),
	}
	err = c.u.ExecTradeV2SwapFromTokenToToken(ctx, trade)
	if err != nil {
		return nil, err
	}

	time.Sleep(3 * time.Second)
	txHash := trade.OrderedTxs[0]
	rx, err := c.UserA.GetTxReceipt(ctx, common.Hash(txHash))
	if err != nil {
		return nil, err
	}
	if rx.Status == types.ReceiptStatusFailed {
		log.Err(err).Interface("token", c.SendContractTxPayload.SmartContractAddr).Msg("tx failed, amount is 0")
		return nil, errors.New("tx failed, amount is 0")
	}
	endBalUserA, err := c.UserA.ReadERC20TokenBalance(ctx, c.SmartContractAddr, c.UserA.Address().String())
	if err != nil {
		return nil, err
	}

	fmt.Println("end balance userA", endBalUserA.String())
	return nil, nil
}
