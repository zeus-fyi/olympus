package web3_client

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (u *UniswapClient) ProcessUniversalRouterTxs(ctx context.Context) {
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
	for methodName, tx := range u.MevSmartContractTxMapUniversalRouter.MethodTxMap {
		switch methodName {
		case V3SwapExactIn:
			u.V3SwapExactIn(tx, tx.Args)
		case V3SwapExactOut:
			u.V3SwapExactOut(tx, tx.Args)
		case V2SwapExactIn:
			u.V2SwapExactIn(tx, tx.Args)
		case V2SwapExactOut:
			u.V2SwapExactOut(tx, tx.Args)
		default:
		}
	}
	fmt.Println("totalFilteredCount:", count)
}
