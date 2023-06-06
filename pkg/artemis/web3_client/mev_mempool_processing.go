package web3_client

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

func DecodeTxArgData(ctx context.Context, tx *types.Transaction, abiFile *abi.ABI, methodFilter *strings_filter.FilterOpts) (string, map[string]interface{}, error) {
	if tx.Data() == nil {
		return "", nil, errors.New("tx data is nil")
	}
	input := tx.Data()
	calldata := input
	if len(calldata) < 4 {
		log.Info().Interface("method", "unknown").Msg("Web3Client| GetFilteredPendingMempoolTxs invalid calldata length")
		return "", nil, errors.New("invalid calldata length")
	}
	sigdata := calldata[:4]
	if abiFile == nil {
		log.Info().Interface("method", "unknown").Msg("Web3Client| GetFilteredPendingMempoolTxs Abi Invalid")
		return "", nil, errors.New("abi invalid")
	}
	method, merr := abiFile.MethodById(sigdata[:4])
	if merr != nil {
		log.Info().Err(merr).Msg("Web3Client| GetFilteredPendingMempoolTxs Method Invalid")
		return "", nil, errors.New("abi method invalid")
	}
	if method == nil {
		log.Info().Interface("method", "unknown").Msg("Web3Client| GetFilteredPendingMempoolTxs Method Invalid")
		return "", nil, errors.New("abi method invalid")
	}
	if method.Name == "" {
		log.Info().Interface("method", "unknown").Msg("Web3Client| GetFilteredPendingMempoolTxs Method Invalid")
		return "", nil, errors.New("abi method name empty")
	}
	if !strings_filter.FilterStringWithOpts(method.Name, methodFilter) {
		//log.Debug().Msg("Web3Client| GetFilteredPendingMempoolTxs Method Filtered")
		return "", nil, errors.New("no methods left after filtering")
	}
	argdata := calldata[4:]
	if len(argdata)%32 != 0 || len(argdata) == 0 {
		//log.Info().Msg("Web3Client| GetFilteredPendingMempoolTxs invalid argdata length")
		return "", nil, errors.New("invalid argdata length")
	}
	m := make(map[string]interface{})
	err := method.Inputs.UnpackIntoMap(m, argdata)
	if err != nil {
		log.Info().Err(err).Msg("Web3Client| UnpackIntoMap invalid")
		return "", nil, errors.New("unpack into map invalid")
	}
	return method.Name, m, nil
}

func ProcessMempoolTxs(ctx context.Context, mempool map[string]map[string]*types.Transaction, mevTxMap MevSmartContractTxMap) (MevSmartContractTxMap, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if mevTxMap.MethodTxMap == nil {
		mevTxMap.MethodTxMap = make(map[string]MevTx)
	}
	var mevTxs []MevTx
	smartContractAddrFilter := common.HexToAddress(mevTxMap.SmartContractAddr)
	smartContractAddrFilterString := smartContractAddrFilter.String()
	for userAddr, txPoolQueue := range mempool {
		for order, tx := range txPoolQueue {
			if tx.To() != nil && tx.To().String() == smartContractAddrFilterString {
				if tx.Data() != nil {
					methodName, args, err := DecodeTxArgData(ctx, tx, mevTxMap.Abi, mevTxMap.Filter)
					if err != nil {
						continue
					}
					singleTx := MevTx{
						UserAddr:    userAddr,
						Args:        args,
						Order:       order,
						TxPoolQueue: txPoolQueue,
						Tx:          tx,
					}
					mevTxs = append(mevTxs, singleTx)
					mevTxMap.MethodTxMap[methodName] = singleTx
				}
			}
		}
	}
	mevTxMap.Txs = append(mevTxMap.Txs, mevTxs...)
	return mevTxMap, nil
}
