package web3_client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
	mempool_txs "github.com/zeus-fyi/olympus/datastores/dynamodb/mev"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func (s *Web3ClientTestSuite) TestRawMempoolTxFilter() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ForceDirToTestDirLocation()
	s.LocalMainnetWeb3User.Web3Actions.Dial()
	defer s.LocalMainnetWeb3User.Close()
	creds := dynamodb_client.DynamoDBCredentials{
		Region:       "us-west-1",
		AccessKey:    s.Tc.AwsAccessKeyDynamoDB,
		AccessSecret: s.Tc.AwsSecretKeyDynamoDB,
	}
	client := mempool_txs.NewMevDynamoDB(creds)
	s.Require().NotNil(client)

	txs, terr := client.GetMempoolTxs(ctx, "mainnet")
	s.Require().Nil(terr)
	s.Require().NotEmpty(txs)
	uni := InitUniswapClient(ctx, s.MainnetWeb3User)
	uni.PrintOn = true
	uni.PrintLocal = true
	uni.Path = filepaths.Path{
		PackageName: "",
		DirIn:       "",
		DirOut:      "./trade_analysis",
		FnIn:        "",
		FnOut:       "",
		Env:         "",
	}
	mempool, err := ConvertMempoolTxs(ctx, txs)
	s.Require().Nil(err)

	err = uni.ProcessMempoolTxs(ctx, mempool)
	s.Require().Nil(err)

	uni.ProcessTxs(ctx)
	count := len(uni.SwapExactTokensForTokensParamsSlice)
	fmt.Println("Total SwapExactTokensForTokensParamsSlice found", len(uni.SwapExactTokensForTokensParamsSlice))
	count += len(uni.SwapTokensForExactTokensParamsSlice)
	fmt.Println("Total SwapTokensForExactTokensParamsSlice found", len(uni.SwapTokensForExactTokensParamsSlice))
	count += len(uni.SwapExactETHForTokensParamsSlice)
	fmt.Println("Total SwapExactETHForTokensParamsSlice found", len(uni.SwapExactETHForTokensParamsSlice))
	count += len(uni.SwapTokensForExactETHParamsSlice)
	fmt.Println("Total SwapTokensForExactETHParamsSlice found", len(uni.SwapTokensForExactETHParamsSlice))
	count += len(uni.SwapExactTokensForETHParamsSlice)
	fmt.Println("Total SwapExactTokensForETHParamsSlice found", len(uni.SwapExactTokensForETHParamsSlice))
	count += len(uni.SwapETHForExactTokensParamsSlice)
	fmt.Println("Total SwapETHForExactTokensParamsSlice found", len(uni.SwapETHForExactTokensParamsSlice))
	fmt.Println("Total trades found", count)
}

func ConvertMempoolTxs(ctx context.Context, mempoolTxs []mempool_txs.MempoolTxsDynamoDB) (map[string]map[string]*types.Transaction, error) {
	txMap := make(map[string]map[string]*types.Transaction)
	for _, tx := range mempoolTxs {
		if txMap[tx.Pubkey] == nil {
			txMap[tx.Pubkey] = make(map[string]*types.Transaction)
		}
		txRpc := &types.Transaction{}
		b, berr := json.Marshal(tx.Tx)
		if berr != nil {
			log.Err(berr).Msg("ConvertMempoolTxs: error marshalling tx")
			return nil, berr
		}
		var txRpcMapStr string
		berr = json.Unmarshal(b, &txRpcMapStr)
		if berr != nil {
			log.Err(berr).Msg("ConvertMempoolTxs: error marshalling tx")
			return nil, berr
		}
		berr = json.Unmarshal([]byte(txRpcMapStr), &txRpc)
		if berr != nil {
			log.Err(berr).Msg("ConvertMempoolTxs: error marshalling tx")
			return nil, berr
		}

		tmp := txMap[tx.Pubkey]
		tmp[fmt.Sprintf("%d", tx.TxOrder)] = txRpc
		txMap[tx.Pubkey] = tmp
	}
	return txMap, nil
}
