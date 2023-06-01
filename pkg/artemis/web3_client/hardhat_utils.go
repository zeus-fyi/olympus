package web3_client

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

func getSlot(userAddress string, slot *big.Int) (string, error) {
	// compute keccak256 hash
	addr := common.HexToAddress(userAddress)
	hash := crypto.Keccak256Hash(
		common.LeftPadBytes(addr.Bytes(), 32),
		common.LeftPadBytes(slot.Bytes(), 32),
	)
	// return hex string of the hash
	return hash.Hex(), nil
}

func (w *Web3Client) GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	w.Dial()
	defer w.Close()
	rx, err := w.C.TransactionReceipt(ctx, txHash)
	if err != nil {
		return rx, err
	}
	return rx, nil
}

func (w *Web3Client) HardhatResetNetworkToBlockBeforeTxMined(ctx context.Context, simNodeUrl string, simNetworkClient, realNetworkClient Web3Client, txHash common.Hash) (int, error) {
	realNetworkClient.Dial()
	rx, err := realNetworkClient.C.TransactionReceipt(ctx, txHash)
	if err != nil {
		return 0, err
	}
	realNetworkClient.Close()
	simNetworkClient.Dial()
	err = simNetworkClient.HardHatResetNetwork(ctx, simNodeUrl, int(rx.BlockNumber.Int64()-1))
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
