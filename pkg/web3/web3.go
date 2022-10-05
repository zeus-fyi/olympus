package web3

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

//curl --user  \
//--header 'Content-Type: application/json' \
//--data '{"jsonrpc":"2.0","method":"eth_getBalance",
//"params":["0xc94770007dda54cF92009BFF0dE90c06F603a09f",
//"latest"],"id":1}' \
//https://868605ce-acde-424e-800c-55ab87808268.ethereum.bison.run

var web3Endpoint = "https://@868605ce-acde-424e-800c-55ab87808268.ethereum.bison.run"

func GetGasPrice(ctx context.Context) (*big.Int, error) {
	cl, err := ethclient.Dial(web3Endpoint)
	if err != nil {
		return nil, err
	}
	gasPrice, err := cl.SuggestGasPrice(ctx)

	if err != nil {
		return nil, err
	}
	return gasPrice, nil
}

func GetTxData(ctx context.Context, txHash string) (*types.Transaction, error) {
	cl, err := ethclient.Dial(web3Endpoint)
	if err != nil {
		return nil, err
	}
	h := common.HexToHash(txHash)
	txData, _, err := cl.TransactionByHash(ctx, h)

	if err != nil {
		return nil, err
	}
	return txData, nil
}
