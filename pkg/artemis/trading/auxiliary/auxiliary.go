package artemis_trading_auxiliary

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/trading/flashbots"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type AuxiliaryTradingUtils struct {
	f          artemis_flashbots.FlashbotsClient
	U          *web3_client.UniswapClient
	Bundle     artemis_flashbots.MevTxBundle
	MevTxGroup MevTxGroup
}

type MevTxGroup struct {
	EventID    int
	OrderedTxs []TxWithMetadata
	MevTxs     []artemis_eth_txs.EthTx
}

type TxWithMetadata struct {
	TradeType       string
	UserOffsetNonce int
	Permit2Tx       artemis_autogen_bases.Permit2Tx
	Tx              *types.Transaction
}

func (m *MevTxGroup) GetRawOrderedTxs() []*types.Transaction {
	txSlice := make([]*types.Transaction, len(m.OrderedTxs))
	for i, tx := range m.OrderedTxs {
		txSlice[i] = tx.Tx
	}
	return txSlice
}

func (m *MevTxGroup) GetHexEncodedTxs() []*types.Transaction {
	txSlice := make([]*types.Transaction, len(m.OrderedTxs))
	for i, tx := range m.OrderedTxs {
		txSlice[i] = tx.Tx
	}
	return txSlice
}

func InitAuxiliaryTradingUtils(ctx context.Context, wa web3_client.Web3Client) AuxiliaryTradingUtils {
	uni := web3_client.InitUniswapClient(ctx, wa)
	fba := artemis_flashbots.InitFlashbotsClient(ctx, &wa.Web3Actions)
	return AuxiliaryTradingUtils{
		U: &uni,
		f: fba,
	}
}

func (a *AuxiliaryTradingUtils) network() string {
	return a.w3c().Network
}

func (a *AuxiliaryTradingUtils) nodeURL() string {
	return a.w3c().NodeURL
}

func (a *AuxiliaryTradingUtils) tradersAccount() *accounts.Account {
	return a.w3c().Account
}

func (a *AuxiliaryTradingUtils) dial() {
	a.w3c().Dial()
}

func (a *AuxiliaryTradingUtils) close() {
	a.w3c().Close()
}

func (a *AuxiliaryTradingUtils) w3c() *web3_client.Web3Client {
	return &a.U.Web3Client
}

func (a *AuxiliaryTradingUtils) w3a() *web3_actions.Web3Actions {
	return &a.w3c().Web3Actions
}

func (a *AuxiliaryTradingUtils) C() *ethclient.Client {
	return a.w3a().C
}

func (a *AuxiliaryTradingUtils) getBlockNumber(ctx context.Context) (int, error) {
	a.dial()
	defer a.close()
	bn, err := a.w3c().C.BlockNumber(ctx)
	if err != nil {
		log.Err(err).Msg("failed to get block number")
		return -1, err
	}
	return int(bn), err
}

const (
	BlockNumber = "BlockNumber"
)

func (a *AuxiliaryTradingUtils) setBlockNumberCtx(ctx context.Context, bn string) context.Context {
	ctx = context.WithValue(ctx, BlockNumber, bn)
	return ctx
}

func (a *AuxiliaryTradingUtils) getBlockNumberCtx(ctx context.Context) string {
	td := ctx.Value(BlockNumber)
	if td != nil {
		return td.(string)
	}
	bn, err := a.getBlockNumber(ctx)
	if err != nil {
		return ""
	}
	bnStr := hexutil.EncodeUint64(uint64(bn + 1))
	return bnStr
}

const (
	TradeDeadline = "TradeDeadline"
)

func (a *AuxiliaryTradingUtils) GetDeadline() *big.Int {
	deadline := int(time.Now().Add(60 * time.Second).Unix())
	sigDeadline := artemis_eth_units.NewBigInt(deadline)
	return sigDeadline
}

func (a *AuxiliaryTradingUtils) getNewTradeDeadlineCtx(ctx context.Context) *big.Int {
	td := ctx.Value(TradeDeadline)
	if td != nil {
		return td.(*big.Int)
	}
	return nil
}

func (a *AuxiliaryTradingUtils) setNewTradeDeadlineCtx(ctx context.Context) context.Context {
	deadline := a.GetDeadline()
	ctx = context.WithValue(ctx, TradeDeadline, deadline)
	return ctx
}
