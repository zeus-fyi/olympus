package web3_client

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (u *UniswapClient) ProcessTxs(ctx context.Context) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Web3Client.Dial()
	bn, err := u.Web3Client.GetHeadBlockHeight(ctx)
	if err != nil {
		log.Err(err).Msg("failed to get block number")
		u.Web3Client.Close()
		return
	}
	u.Web3Client.Close()
	u.BlockNumber = bn
	count := 0
	for _, tx := range u.MevSmartContractTxMapUniversalRouter.Txs {
		u.ProcessUniversalRouterTxs(ctx, tx)
	}
	for _, tx := range u.MevSmartContractTxMap.Txs {
		switch tx.MethodName {
		case addLiquidity:
			//u.AddLiquidity(tx.Args)
		case addLiquidityETH:
			// payable
			//u.AddLiquidityETH(tx.Args)
			if tx.Tx.Value() == nil {
				continue
			}
		case removeLiquidity:
			//u.RemoveLiquidity(tx.Args)
		case removeLiquidityETH:
			//u.RemoveLiquidityETH(tx.Args)
		case removeLiquidityWithPermit:
			//u.RemoveLiquidityWithPermit(tx.Args)
		case removeLiquidityETHWithPermit:
			//u.RemoveLiquidityETHWithPermit(tx.Args)
		case swapExactTokensForTokens:
			count++
			u.SwapExactTokensForTokens(tx, tx.Args)
		case swapTokensForExactTokens:
			count++
			u.SwapTokensForExactTokens(tx, tx.Args)
		case swapExactETHForTokens:
			// payable
			count++
			if tx.Tx.Value() == nil {
				continue
			}
			u.SwapExactETHForTokens(tx, tx.Args, tx.Tx.Value())
		case swapTokensForExactETH:
			count++
			u.SwapTokensForExactETH(tx, tx.Args)
		case swapExactTokensForETH:
			count++
			u.SwapExactTokensForETH(tx, tx.Args)
		case swapETHForExactTokens:
			// payable
			count++
			if tx.Tx.Value() == nil {
				continue
			}
			u.SwapETHForExactTokens(tx, tx.Args, tx.Tx.Value())
		}
	}
	fmt.Println("totalFilteredCount:", count)
}
