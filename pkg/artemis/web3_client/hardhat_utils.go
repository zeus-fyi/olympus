package web3_client

import (
	"context"
	"math/big"

	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/gochain/v4/common/hexutil"
	"github.com/gochain/gochain/v4/crypto"
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

func (w *Web3Client) HardhatResetNetworkToBlockBeforeTxMined(ctx context.Context, simNodeUrl string, simNetworkClient, realNetworkClient Web3Client, txHash common.Hash) error {
	realNetworkClient.Dial()
	rx, err := realNetworkClient.Client.GetTransactionByHash(ctx, txHash)
	if err != nil {
		return err
	}
	realNetworkClient.Close()
	simNetworkClient.Dial()
	err = simNetworkClient.Client.ResetNetwork(ctx, simNodeUrl, int(rx.BlockNumber.Int64()-1))
	if err != nil {
		return err
	}
	simNetworkClient.Close()
	return nil
}
