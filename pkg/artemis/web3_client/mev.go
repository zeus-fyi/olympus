package web3_client

import (
	"context"

	"github.com/gochain/gochain/v4/accounts/abi"
	"github.com/gochain/gochain/v4/common"
	"github.com/rs/zerolog/log"
	web3_types "github.com/zeus-fyi/gochain/web3/types"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

type MevSmartContractTxMap struct {
	SmartContractAddr string
	Abi               *abi.ABI
	MethodTxMap       map[string]MevTx
	Txs               []MevTx
	Filter            *strings_filter.FilterOpts
}

type MevTx struct {
	UserAddr    string
	Args        map[string]interface{}
	Order       string
	TxPoolQueue map[string]*web3_types.RpcTransaction
	Tx          *web3_types.RpcTransaction
}

func (w *Web3Client) GetFilteredPendingMempoolTxs(ctx context.Context, mevTxMap MevSmartContractTxMap) (MevSmartContractTxMap, error) {
	if mevTxMap.MethodTxMap == nil {
		mevTxMap.MethodTxMap = make(map[string]MevTx)
	}
	w.Web3Actions.Dial()
	defer w.Close()
	mempool, err := w.Web3Actions.GetTxPoolContent(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Web3Client| GetFilteredPendingMempoolTxs")
		return mevTxMap, err
	}
	var mevTxs []MevTx
	smartContractAddrFilter := common.HexToAddress(mevTxMap.SmartContractAddr)
	smartContractAddrFilterString := smartContractAddrFilter.String()
	for userAddr, txPoolQueue := range mempool["pending"] {
		for order, tx := range txPoolQueue {
			if tx.To != nil && tx.To.String() == smartContractAddrFilterString {
				if tx.Input != nil {
					input := *tx.Input
					calldata := []byte(input)
					if len(calldata) < 4 {
						log.Info().Interface("method", "unknown").Msg("Web3Client| GetFilteredPendingMempoolTxs invalid calldata length")
						continue
					}
					sigdata := calldata[:4]
					method, merr := mevTxMap.Abi.MethodById(sigdata[:4])
					if merr != nil {
						log.Info().Err(err).Interface("method", method.Name).Msg("Web3Client| GetFilteredPendingMempoolTxs Method Invalid")
						continue
					}
					if !strings_filter.FilterStringWithOpts(method.Name, mevTxMap.Filter) {
						log.Info().Interface("method", method.Name).Msg("Web3Client| GetFilteredPendingMempoolTxs Method Filtered")
						continue
					}
					argdata := calldata[4:]
					if len(argdata)%32 != 0 {
						log.Info().Interface("method", method.Name).Msg("Web3Client| GetFilteredPendingMempoolTxs invalid argdata length")
						continue
					}
					m := make(map[string]interface{})
					err = method.Inputs.UnpackIntoMap(m, argdata)
					if err != nil {
						log.Info().Err(err).Interface("method", method.Name).Msg("Web3Client| UnpackIntoMap invalid")
						continue
					}
					singleTx := MevTx{
						UserAddr:    userAddr,
						Args:        m,
						Order:       order,
						TxPoolQueue: txPoolQueue,
						Tx:          tx,
					}
					mevTxs = append(mevTxs, singleTx)
					mevTxMap.MethodTxMap[method.Name] = singleTx
				}
			}
		}
	}
	mevTxMap.Txs = mevTxs
	return mevTxMap, nil
}
