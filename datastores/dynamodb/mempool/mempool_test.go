package mempool_txs

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/gochain/v4/common/hexutil"
	"github.com/stretchr/testify/suite"
	web3_types "github.com/zeus-fyi/gochain/web3/types"
	"github.com/zeus-fyi/olympus/configs"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type MempoolTxDynamoDBTestSuite struct {
	m MempoolTxsDynamoDB
	test_suites_base.TestSuite
}

var (
	ctx    = context.Background()
	region = "us-west-1"
)

func (t *MempoolTxDynamoDBTestSuite) SetupTest() {
	t.Tc = configs.InitLocalTestConfigs()
}

func (t *MempoolTxDynamoDBTestSuite) TestGetMempoolTxs() {
	creds := dynamodb_client.DynamoDBCredentials{
		Region:       region,
		AccessKey:    t.Tc.AwsAccessKeyDynamoDB,
		AccessSecret: t.Tc.AwsSecretKeyDynamoDB,
	}
	m := NewMempoolTxDynamoDB(creds)
	memTxs, err := m.GetMempoolTxs(ctx, "mainnet")
	t.Require().Nil(err)
	t.Require().NotNil(memTxs)

	txMap := make(map[string]map[string]*web3_types.RpcTransaction)

	for _, tx := range memTxs {
		if txMap[tx.Pubkey] == nil {
			txMap[tx.Pubkey] = make(map[string]*web3_types.RpcTransaction)
		}
		var txRpcMap string
		txRpc := &Transaction{}
		b, berr := json.Marshal(tx.Tx)
		t.Require().Nil(berr)
		berr = json.Unmarshal(b, &txRpcMap)
		t.Require().Nil(berr)
		t.Require().NotNil(txRpcMap)
		berr = json.Unmarshal([]byte(txRpcMap), &txRpc)
		t.Require().Nil(berr)
		txGoRpc, berr := ProcessTransaction(*txRpc)
		t.Require().Nil(berr)
		tmp := txMap[tx.Pubkey]
		tmp[fmt.Sprintf("%d", tx.TxOrder)] = txGoRpc
		txMap[tx.Pubkey] = tmp
	}

	t.Assert().NotEmpty(txMap)
}

func TestMempoolTxDynamoDBTestSuite(t *testing.T) {
	suite.Run(t, new(MempoolTxDynamoDBTestSuite))
}

type Transaction struct {
	Type                 string        `json:"type"`
	Nonce                string        `json:"nonce"`
	GasPrice             interface{}   `json:"gasPrice"`
	MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         string        `json:"maxFeePerGas"`
	Gas                  string        `json:"gas"`
	Value                string        `json:"value"`
	Input                string        `json:"input"`
	V                    string        `json:"v"`
	R                    string        `json:"r"`
	S                    string        `json:"s"`
	To                   string        `json:"to"`
	ChainID              string        `json:"chainId"`
	AccessList           []interface{} `json:"accessList"`
	Hash                 string        `json:"hash"`
}

func ProcessTransaction(tx Transaction) (*web3_types.RpcTransaction, error) {
	rpcTx := &web3_types.RpcTransaction{}
	nonce, err := hexutil.DecodeUint64(tx.Nonce)
	if err != nil {
		return nil, err
	}
	rpcTx.Nonce = (*hexutil.Uint64)(&nonce)

	gasLimit, err := hexutil.DecodeUint64(tx.Gas)
	if err != nil {
		return nil, err
	}
	rpcTx.GasLimit = (*hexutil.Uint64)(&gasLimit)
	to := common.HexToAddress(tx.To)
	rpcTx.To = &to
	if tx.Input != "" {
		input := hexutil.Bytes(common.Hex2Bytes(tx.Input))
		rpcTx.Input = &input
	}
	if tx.Hash != "" {
		hash := common.HexToHash(tx.Hash)
		rpcTx.Hash = &hash
	}
	return rpcTx, nil
}

/*
// CopyFrom copies the fields from t to r.
func (r *RpcTransaction) CopyFrom(t *Transaction) {
	r.Nonce = (*hexutil.Uint64)(&t.Nonce)
	r.GasPrice = (*hexutil.Big)(t.GasPrice)
	r.GasLimit = (*hexutil.Uint64)(&t.GasLimit)
	r.To = t.To
	r.Value = (*hexutil.Big)(t.Value)
	r.Input = (*hexutil.Bytes)(&t.Input)
	r.Hash = &t.Hash
	r.BlockNumber = (*hexutil.Big)(t.BlockNumber)
	r.BlockHash = &t.BlockHash
	r.From = &t.From
	r.TransactionIndex = (*hexutil.Uint64)(&t.TransactionIndex)
	r.V = (*hexutil.Big)(t.V)
	r.R = (*hexutil.Big)(t.R)
	r.S = (*hexutil.Big)(t.S)
}

*/
