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
	for _, tx := range u.MevSmartContractTxMapUniversalRouterOld.Txs {
		u.ProcessUniversalRouterTxs(ctx, tx)
	}
	for _, tx := range u.MevSmartContractTxMapUniversalRouterNew.Txs {
		u.ProcessUniversalRouterTxs(ctx, tx)
	}
	for _, tx := range u.MevSmartContractTxMapV3SwapRouterV1.Txs {
		u.ProcessUniswapV3RouterTxs(ctx, tx)
	}
	for _, tx := range u.MevSmartContractTxMapV3SwapRouterV2.Txs {
		u.ProcessUniswapV3RouterTxs(ctx, tx)
	}
	u.ProcessV2Router01Txs()
	u.ProcessV2Router02Txs()
	fmt.Println("totalFilteredCount:", count)
}
