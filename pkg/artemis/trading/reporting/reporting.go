package artemis_reporting

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	uniswap_core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
)

/*
SELECT count(*), amount_out_addr, expected_profit_amount_out
FROM eth_mev_tx_analysis
WHERE end_reason = 'success' AND rx_block_number > 17639300 AND amount_in_addr = '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2'
GROUP BY amount_out_addr, expected_profit_amount_out
*/
type RewardsGroup struct {
	Map map[string]RewardsHistory
}

type RewardsHistory struct {
	Count                   int
	FailedCount             int
	AmountOutToken          *uniswap_core_entities.Token
	ExpectedProfitAmountOut *big.Int
}

func getQ(blockCheckpoint int, tradeMethod string) string {

	var add string
	if tradeMethod == "any" {
	} else {
		add = fmt.Sprintf("AND trade_method = '%s'", tradeMethod)
	}
	var que = fmt.Sprintf(`WITH cte_failed_count AS (
			SELECT COUNT(*) as failed, amount_out_addr
			FROM eth_mev_tx_analysis
			WHERE end_reason != 'success'
		    AND (end_reason = 'expected minus actual profit mismatch' OR end_reason = 'execution reverted' OR end_reason = 'execution reverted: TRANSFER_FROM_FAILED' OR end_reason = 'execution reverted: UniswapV2: TRANSFER_FAILED') 
			AND rx_block_number > %d AND amount_in_addr = '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' %s
			GROUP BY amount_out_addr
			) 
				SELECT me.amount_out_addr, expected_profit_amount_out, et.name, et.symbol, et.decimals, et.transfer_tax_numerator, et.transfer_tax_denominator, COALESCE(fc.failed, 0) as failed_count
				FROM eth_mev_tx_analysis me
				INNER JOIN erc20_token_info et ON et.address = me.amount_out_addr
				LEFT JOIN cte_failed_count fc ON fc.amount_out_addr = me.amount_out_addr
				WHERE end_reason = 'success' AND rx_block_number > %d AND expected_profit_amount_out IS NOT NULL AND amount_in_addr = '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' %s
			`, blockCheckpoint, add, blockCheckpoint, add)
	return que
}

type RewardHistoryFilter struct {
	TradeMethod string
	FromBlock   int
}

func GetRewardsHistory(ctx context.Context, rhf RewardHistoryFilter) (*RewardsGroup, error) {
	rw := &RewardsGroup{
		Map: make(map[string]RewardsHistory),
	}
	q := getQ(rhf.FromBlock, rhf.TradeMethod)
	rows, err := apps.Pg.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		rh := RewardsHistory{}
		var addrOut sql.NullString
		var profit sql.NullString
		var name sql.NullString
		var symbol sql.NullString
		var decimals sql.NullInt64
		var transferTaxNumerator sql.NullInt64
		var transferTaxDenominator sql.NullInt64

		rowErr := rows.Scan(&addrOut, &profit, &name, &symbol, &decimals, &transferTaxNumerator, &transferTaxDenominator, &rh.FailedCount)
		if rowErr != nil {
			log.Err(rowErr).Msg("GetRewardsHistory")
			return nil, rowErr
		}

		if transferTaxNumerator.Valid && transferTaxDenominator.Valid {
			tt := uniswap_core_entities.NewFraction(artemis_eth_units.NewBigInt(int(transferTaxNumerator.Int64)), artemis_eth_units.NewBigInt(int(transferTaxDenominator.Int64)))

			if name.Valid == false || symbol.Valid == false || decimals.Valid == false || profit.Valid == false || addrOut.Valid == false {
				//fmt.Println("name, symbol, decimals are nil", addrOut)
				continue
			}
			rh.AmountOutToken = uniswap_core_entities.NewTokenWithTransferTax(1, accounts.HexToAddress(addrOut.String), uint(decimals.Int64), symbol.String, name.String, tt)
			rh.ExpectedProfitAmountOut = artemis_eth_units.NewBigIntFromStr(profit.String)
			if v, ok := rw.Map[addrOut.String]; ok {
				v.Count += 1
				v.ExpectedProfitAmountOut = artemis_eth_units.AddBigInt(v.ExpectedProfitAmountOut, rh.ExpectedProfitAmountOut)
				rw.Map[addrOut.String] = v
			} else {
				rh.Count = 1
				rw.Map[addrOut.String] = rh
			}
		}
	}
	return rw, nil
}
