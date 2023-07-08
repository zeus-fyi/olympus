package artemis_reporting

import (
	"context"
	"math/big"

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
	AmountOutToken          uniswap_core_entities.Token
	ExpectedProfitAmountOut *big.Int
}

var q = `SELECT amount_out_addr, expected_profit_amount_out
		FROM eth_mev_tx_analysis
		WHERE end_reason = 'success' AND rx_block_number > 17639300 AND amount_in_addr = '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2'`

func GetRewardsHistory(ctx context.Context) (*RewardsGroup, error) {
	rw := &RewardsGroup{
		Map: make(map[string]RewardsHistory),
	}
	rows, err := apps.Pg.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		rh := RewardsHistory{}
		addrOut := ""
		profit := ""
		rowErr := rows.Scan(&addrOut, &profit)
		if rowErr != nil {
			return nil, rowErr
		}
		rh.AmountOutToken.Address = accounts.HexToAddress(addrOut)
		rh.ExpectedProfitAmountOut = artemis_eth_units.NewBigIntFromStr(profit)
		if v, ok := rw.Map[addrOut]; ok {
			v.Count += 1
			v.ExpectedProfitAmountOut = artemis_eth_units.AddBigInt(v.ExpectedProfitAmountOut, rh.ExpectedProfitAmountOut)
			rw.Map[addrOut] = v
		} else {
			rw.Map[addrOut] = rh
		}
	}
	return rw, nil
}
