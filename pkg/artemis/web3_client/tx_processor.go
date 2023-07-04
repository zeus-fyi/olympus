package web3_client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
)

func (u *UniswapClient) ProcessTxs(ctx context.Context) {
	wc := web3_actions.NewWeb3ActionsClient(artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLive.NodeURL)
	wc.Dial()
	bn, berr := wc.C.BlockNumber(ctx)
	if berr != nil {
		log.Err(berr)
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	u.BlockNumber = new(big.Int).SetUint64(bn)
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
