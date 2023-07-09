package artemis_trading_auxiliary

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

var TradingAuxiliary AuxiliaryTradingUtils

type AuxiliaryTradingUtils struct {
	web3_client.Web3Client
	OrderedTxs []accounts.Hash
}

func InitAuxiliaryTradingUtils(ctx context.Context, nodeURL, network string, acc accounts.Account) AuxiliaryTradingUtils {
	TradingAuxiliary = AuxiliaryTradingUtils{
		Web3Client: web3_client.NewWeb3Client(nodeURL, &acc),
	}
	TradingAuxiliary.Network = network
	return TradingAuxiliary
}

func (a *AuxiliaryTradingUtils) AddTx(tx *types.Transaction) {
	if a.OrderedTxs == nil {
		a.OrderedTxs = []accounts.Hash{}
	}
	a.OrderedTxs = append(a.OrderedTxs, accounts.HexToHash(tx.Hash().Hex()))
	a.trackTx(tx)
}
