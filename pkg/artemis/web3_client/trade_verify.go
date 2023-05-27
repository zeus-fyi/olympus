package web3_client

import (
	"fmt"
)

func (u *UniswapV2Client) VerifyTradeResults(tf *TradeExecutionFlowInBigInt) error {
	if u.DebugPrint {
		fmt.Println("executing full sandwich trade")
	}

	return nil
}
