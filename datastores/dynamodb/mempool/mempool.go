package mempool_txs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type MempoolTxDynamoDB struct {
	*dynamodb.Client
	*web3_client.Web3Client
}

func NewMempoolTxDynamoDB(creds dynamodb_client.DynamoDBCredentials, wc *web3_client.Web3Client) MempoolTxDynamoDB {
	d, err := dynamodb_client.NewDynamoDBClient(context.Background(), creds)
	if err != nil {
		log.Err(err)
	}
	return MempoolTxDynamoDB{
		d.Client, wc,
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

func (m *MempoolTxDynamoDB) GetMempoolTxs(ctx context.Context, network string) (map[string]map[string]*types.Transaction, error) {
	var mempoolTxsTableName *string
	if network == "mainnet" {
		mempoolTxsTableName = MainnetMempoolTxsTableName
	} else if network == "goerli" {
		mempoolTxsTableName = GoerliMempoolTxsTableName
	}
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

	txMap := make(map[string]map[string]*types.Transaction)
	for _, tx := range mempoolTxs {
		if txMap[tx.Pubkey] == nil {
			txMap[tx.Pubkey] = make(map[string]*types.Transaction)
		}
		txRpc := &types.Transaction{}
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
		isPending, perr := m.ValidateTxIsPending(ctx, txRpc.Hash().Hex())
		if perr != nil {
			log.Err(perr).Msg("GetMempoolTxs: error validating tx")
			continue
		}
		if !isPending {
			log.Warn().Msg("GetMempoolTxs: tx not pending")
			continue
		}
		tmp := txMap[tx.Pubkey]
		tmp[fmt.Sprintf("%d", tx.TxOrder)] = txRpc
		txMap[tx.Pubkey] = tmp
	}
	//fmt.Println("startingTxCount", len(mempoolTxs))
	//fmt.Println("endFilteredTxCount", len(txMap))
	//fmt.Println("filteredCount", len(mempoolTxs)-len(txMap))
	return txMap, nil
}
