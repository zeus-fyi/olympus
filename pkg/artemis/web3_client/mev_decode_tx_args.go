package web3_client

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

func DecodeTxArgData(ctx context.Context, tx *types.Transaction, mevMap MevSmartContractTxMap) (string, map[string]interface{}, error) {
	if tx.Data() == nil {
		return "", nil, errors.New("tx data is nil")
	}
	input := tx.Data()
	return DecodeTxData(ctx, input, mevMap)
}

func DecodeTxData(ctx context.Context, input []byte, mevMap MevSmartContractTxMap) (string, map[string]interface{}, error) {
	calldata := input
	if len(calldata) < 4 {
		log.Info().Interface("method", "unknown").Msg("Web3Client| GetFilteredPendingMempoolTxs invalid calldata length")
		return "", nil, errors.New("invalid calldata length")
	}
	sigdata := calldata[:4]
	if mevMap.Abi == nil {
		log.Info().Interface("method", "unknown").Msg("Web3Client| GetFilteredPendingMempoolTxs Abi Invalid")
		return "", nil, errors.New("abi invalid")
	}
	method, merr := mevMap.Abi.MethodById(sigdata[:4])
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
	if !strings_filter.FilterStringWithOpts(method.Name, mevMap.Filter) {
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
