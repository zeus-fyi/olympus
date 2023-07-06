package artemis_trade_debugger

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/pkg/artemis/trading/async_analysis"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

func (t *TradeDebugger) analyzeToken(ctx context.Context, address accounts.Address, amountTraded *big.Int) error {
	token := address.String()
	if _, ok := artemis_trading_cache.TokenMap[token]; !ok {
		return errors.New("token not found")
	}

	info := artemis_trading_cache.TokenMap[token]
	if info.Name != nil {
		fmt.Println("ANALYZING token: ", token, "name: ", *info.Name)
	}
	den := info.TransferTaxDenominator
	num := info.TransferTaxNumerator
	fmt.Println("token: ", token, "tradingTax: ", den, "num: ", num)

	ca := async_analysis.NewContractAnalysis(t.UniswapClient, address.String(), artemis_oly_contract_abis.MustLoadERC20Abi())
	ca.UserA = t.UniswapClient.Web3Client
	ca.UserB = t.LiveNetworkClient
	//  -1 means balanceOf it wasn't cracked within 100 slots, -2 means cracking hasn't been attempted yet
	if info.BalanceOfSlotNum == -2 {
		err := ca.FindERC20BalanceOfSlotNumber(ctx)
		if err != nil {
			return err
		}
	}
	if num == nil || den == nil {
		feePerc, err := ca.CalculateTransferFeeTax(ctx, amountTraded)
		if err != nil {
			return err
		}
		fmt.Println("feePerc: ", feePerc.Numerator, feePerc.Denominator)
		err = ca.CalculateTransferFeeTaxRange(ctx)
		if err != nil {
			return err
		}
	}
	t.ContractAnalysis = ca
	return nil
}
