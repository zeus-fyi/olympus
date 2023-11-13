package web3_client

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

const (
	Owner = "owner"
)

func (w *Web3Client) GetOwner(ctx context.Context, abiFile *abi.ABI, contractAddress string) (common.Address, error) {
	w.Dial()
	defer w.C.Close()
	payload := web3_actions.SendContractTxPayload{
		SmartContractAddr: contractAddress,
		ContractABI:       abiFile,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		MethodName:        Owner,
	}
	payload.Params = []interface{}{}
	owner, err := w.GetContractConst(ctx, &payload)
	if err != nil {
		log.Err(err).Msg("GetOwner")
		return common.Address{}, err
	}
	return owner[0].(common.Address), err
}

func (w *Web3Client) EthCall(ctx context.Context, from common.Address, payload *web3_actions.SendContractTxPayload, bn *big.Int) ([]byte, error) {
	w.Dial()
	defer w.Close()
	if payload.Data == nil {
		payload.Data = []byte{}
		err := payload.GenerateBinDataFromParamsAbi(ctx)
		if err != nil {
			log.Err(err).Msg("EthCall: GenerateBinDataFromParamsAbi")
			return nil, err
		}
	}
	toAddr := common.HexToAddress(payload.SmartContractAddr)
	msg := ethereum.CallMsg{
		From:       from,
		To:         &toAddr,
		Gas:        payload.GasLimit,
		GasPrice:   payload.GasPrice,
		GasFeeCap:  payload.GasFeeCap,
		GasTipCap:  payload.GasTipCap,
		Value:      payload.SendEtherPayload.Amount,
		Data:       payload.Data,
		AccessList: nil,
	}
	resp, err := w.C.CallContract(ctx, msg, bn)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("EthCall")
		return nil, err
	}
	return resp, nil
}
