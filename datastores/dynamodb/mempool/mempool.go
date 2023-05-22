package mempool_txs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/gochain/v4/common/hexutil"
	"github.com/rs/zerolog/log"
	web3_types "github.com/zeus-fyi/gochain/web3/types"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
)

type MempoolTxDynamoDB struct {
	*dynamodb.Client
}

func NewMempoolTxDynamoDB(creds dynamodb_client.DynamoDBCredentials) MempoolTxDynamoDB {
	d, err := dynamodb_client.NewDynamoDBClient(context.Background(), creds)
	if err != nil {
		log.Err(err)
	}
	return MempoolTxDynamoDB{
		d.Client,
	}
}

var (
	MainnetMempoolTxsTableName = aws.String("MempoolTxsMainnet")
	GoerliMempoolTxsTableName  = aws.String("MempoolTxsGoerli")
)

type MempoolTxDynamoDBTableKeys struct {
	Pubkey  string `dynamodbav:"pubkey"`
	TxOrder int    `dynamodbav:"txOrder"`
}

type MempoolTxsDynamoDB struct {
	MempoolTxDynamoDBTableKeys
	Tx  string `dynamodbav:"tx"`
	TTL int    `dynamodbav:"ttl"`
}

func (m *MempoolTxDynamoDB) GetMempoolTxs(ctx context.Context, network string) (map[string]map[string]*web3_types.RpcTransaction, error) {
	var mempoolTxsTableName *string
	if network == "mainnet" {
		mempoolTxsTableName = MainnetMempoolTxsTableName
	} else if network == "goerli" {
		mempoolTxsTableName = GoerliMempoolTxsTableName
	}
	fmt.Println(*mempoolTxsTableName)
	scanInput := &dynamodb.ScanInput{
		TableName: mempoolTxsTableName,
	}
	r, err := m.Scan(ctx, scanInput)
	if err != nil {
		log.Err(err).Msg("GetMempoolTxs: error scanning mempool txs")
		return nil, err
	}

	var mempoolTxs []MempoolTxsDynamoDB
	err = attributevalue.UnmarshalListOfMaps(r.Items, &mempoolTxs)
	if err != nil {
		log.Err(err).Msg("GetMempoolTxs: error UnmarshalListOfMaps mempool txs")
		return nil, err
	}

	txMap := make(map[string]map[string]*web3_types.RpcTransaction)
	for _, tx := range mempoolTxs {
		if txMap[tx.Pubkey] == nil {
			txMap[tx.Pubkey] = make(map[string]*web3_types.RpcTransaction)
		}
		txRpc := &Transaction{}
		b, berr := json.Marshal(tx.Tx)
		if berr != nil {
			log.Err(berr).Msg("GetMempoolTxs: error marshalling tx")
			return nil, berr
		}
		var txRpcMapStr string
		berr = json.Unmarshal(b, &txRpcMapStr)
		if berr != nil {
			log.Err(berr).Msg("GetMempoolTxs: error marshalling tx")
			return nil, berr
		}
		berr = json.Unmarshal([]byte(txRpcMapStr), &txRpc)
		if berr != nil {
			log.Err(berr).Msg("GetMempoolTxs: error marshalling tx")
			return nil, berr
		}
		berr = json.Unmarshal(b, &txRpcMapStr)
		if berr != nil {
			log.Err(berr).Msg("GetMempoolTxs: error marshalling tx")
			return nil, berr
		}
		txGoRpc, berr := ProcessTransaction(*txRpc)
		if berr != nil {
			log.Err(berr).Msg("GetMempoolTxs: error marshalling tx")
			return nil, berr
		}

		tmp := txMap[tx.Pubkey]
		from := common.HexToAddress(tx.Pubkey)
		txGoRpc.From = &from
		tmp[fmt.Sprintf("%d", tx.TxOrder)] = txGoRpc
		txMap[tx.Pubkey] = tmp
	}

	return txMap, nil
}

type Transaction struct {
	Type                 string          `json:"type"`
	Nonce                *hexutil.Uint64 `json:"nonce"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas,omitempty"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas,omitempty"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	Value                *hexutil.Big    `json:"value"`
	Input                *hexutil.Bytes  `json:"input"`
	V                    *hexutil.Big    `json:"v"`
	R                    *hexutil.Big    `json:"r"`
	S                    *hexutil.Big    `json:"s"`
	To                   *common.Address `json:"to"`
	ChainID              string          `json:"chainId"`
	AccessList           []interface{}   `json:"accessList"`
	Hash                 *common.Hash    `json:"hash"`
}

func ProcessTransaction(tx Transaction) (*web3_types.RpcTransaction, error) {
	rpcTx := &web3_types.RpcTransaction{
		Nonce:     tx.Nonce,
		GasPrice:  tx.GasPrice,
		GasLimit:  tx.Gas,
		GasFeeCap: tx.MaxFeePerGas,
		GasTipCap: tx.MaxPriorityFeePerGas,
		To:        tx.To,
		Value:     tx.Value,
		Input:     tx.Input,
		V:         tx.V,
		R:         tx.R,
		S:         tx.S,
		Hash:      tx.Hash,
	}
	return rpcTx, nil
}
