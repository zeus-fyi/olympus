package web3_client

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (u *UniswapV2Client) ProcessUniversalRouterTxs(ctx context.Context) {
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
			fmt.Println(tx)
		case V3SwapExactOut:
		case V2SwapExactIn:
		case V2SwapExactOut:
		case Permit2TransferFrom:
		case Permit2PermitBatch:
		case Permit2TransferFromBatch:
		default:
		}

	}
	fmt.Println("totalFilteredCount:", count)
}
