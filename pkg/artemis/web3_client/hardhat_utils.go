package web3_client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (w *Web3Client) MineNextBlock(ctx context.Context) error {
	w.Dial()
	defer w.Close()
	oneBlock := hexutil.Big{}
	bigInt := oneBlock.ToInt()
	bigInt.Set(new(big.Int).SetUint64(0))
	oneBlock = hexutil.Big(*bigInt)
	err := w.MineBlock(ctx, oneBlock)
	if err != nil {
		return err
	}
	return nil
}

func GetSlot(userAddress string, slot *big.Int) (string, error) {
	// compute keccak256 hash
	addr := common.HexToAddress(userAddress)
	hash := crypto.Keccak256Hash(
		common.LeftPadBytes(addr.Bytes(), 32),
		common.LeftPadBytes(slot.Bytes(), 32),
	)
	// return hex string of the hash
	return hash.Hex(), nil
}

func (w *Web3Client) PendingNonce(ctx context.Context, user common.Address) (int, error) {
	w.Dial()
	defer w.Close()
	nonce, err := w.C.PendingNonceAt(ctx, user)
	if err != nil {
		return -1, err
	}
	return int(nonce), nil
}

func (w *Web3Client) NonceAt(ctx context.Context, user common.Address, bn *big.Int) (int, error) {
	w.Dial()
	defer w.Close()
	nonce, err := w.C.NonceAt(ctx, user, bn)
	if err != nil {
		return -1, err
	}
	return int(nonce), nil
}

func (w *Web3Client) HardhatResetNetworkToBlock(ctx context.Context, blockNum int) error {
	w.Dial()
	defer w.Close()
	err := w.HardHatResetNetwork(ctx, blockNum)
	if err != nil {
		return err
	}
	return nil
}
func (w *Web3Client) HardhatResetNetworkToBlockBeforeTxMined(ctx context.Context, simNodeUrl string, simNetworkClient, realNetworkClient Web3Client, txHash common.Hash) (int, error) {
	realNetworkClient.Dial()
	rx, err := realNetworkClient.C.TransactionReceipt(ctx, txHash)
	if err != nil {
		return 0, err
	}
	realNetworkClient.Close()
	simNetworkClient.Dial()
	err = simNetworkClient.HardHatResetNetwork(ctx, int(rx.BlockNumber.Int64()-1))
	if err != nil {
		return 0, err
	}
	simNetworkClient.Close()
	return int(rx.BlockNumber.Int64()), nil
}

func (w *Web3Client) SendImpersonatedTx(ctx context.Context, tx *types.Transaction) error {
	sender := types.LatestSignerForChainID(tx.ChainId())
	from, err := sender.Sender(tx)
	if err != nil {
		return err
	}
	err = w.HardhatImpersonateAccount(ctx, from.String())
	if err != nil {
		return err
	}
	err = w.SendSignedTransaction(ctx, tx)
	if err != nil {
		return err
	}
	err = w.StopImpersonatingAccount(ctx, from.String())
	if err != nil {
		return err
	}
	return nil
}

func (w *Web3Client) GetBlockTxs(ctx context.Context) (types.Transactions, error) {
	w.Dial()
	defer w.Close()
	block, err := w.C.BlockByNumber(ctx, nil)
	if err != nil {
		log.Err(err).Msg("failed to get block txs")
		return nil, err
	}
	return block.Transactions(), nil
}

func (w *Web3Client) GetTxByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	w.Dial()
	defer w.Close()
	tx, isPending, err := w.C.TransactionByHash(ctx, hash)
	if err != nil {
		log.Err(err).Msg("failed to get tx by hash")
		return nil, false, err
	}
	return tx, isPending, nil
}

func (w *Web3Client) GetBlockHeight(ctx context.Context) (*big.Int, error) {
	w.Dial()
	defer w.Close()
	blockNumber, err := w.C.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetUint64(blockNumber), nil
}

func (w *Web3Client) GetNodeMetadata(ctx context.Context) (web3_actions.NodeInfo, error) {
	w.Dial()
	defer w.Close()
	info, err := w.GetNodeInfo(ctx)
	if err != nil {
		log.Err(err).Msg("failed to get node info")
		return info, err
	}
	return info, nil
}

func (w *Web3Client) GetHeadBlockHeight(ctx context.Context) (*big.Int, error) {
	w.Dial()
	defer w.C.Close()
	isSyncing, err := w.C.SyncProgress(ctx)
	if err != nil {
		fmt.Println("error getting sync status", isSyncing)
		return nil, err
	}
	if isSyncing != nil && isSyncing.CurrentBlock != isSyncing.HighestBlock {
		return nil, fmt.Errorf("node is not synced")
	}
	return w.GetBlockHeight(ctx)
}
