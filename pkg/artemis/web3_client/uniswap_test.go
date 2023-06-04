package web3_client

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func (s *Web3ClientTestSuite) TestUniswapMempoolFilter() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ForceDirToTestDirLocation()
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
	txMap, err := s.MainnetWeb3User.GetFilteredPendingMempoolTxs(ctx, uni.MevSmartContractTxMap)
	s.Require().Nil(err)
	s.Assert().NotEmpty(txMap)
	uni.MevSmartContractTxMap = txMap
	uni.ProcessTxs(ctx)
	count := len(uni.SwapExactTokensForTokensParamsSlice)
	count += len(uni.SwapTokensForExactTokensParamsSlice)
	count += len(uni.SwapExactETHForTokensParamsSlice)
	count += len(uni.SwapTokensForExactETHParamsSlice)
	count += len(uni.SwapExactTokensForETHParamsSlice)
	count += len(uni.SwapETHForExactTokensParamsSlice)
	fmt.Println("Total trades found", count)
}

func (s *Web3ClientTestSuite) TestMevTxSelect() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	mevTxs, err := artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17275807)
	s.Require().Nil(err)
	s.Require().NotEmpty(mevTxs)

	for _, mevTx := range mevTxs {
		tf := TradeExecutionFlow{}
		b := []byte(mevTx.TxFlowPrediction)
		berr := json.Unmarshal(b, &tf)
		s.Require().Nil(berr)
		s.Require().NotEmpty(tf.UserTrade)

		txRpc := types.Transaction{}
		b = []byte(mevTx.Tx)
		berr = json.Unmarshal(b, &txRpc)
		s.Require().Nil(berr)
		s.Require().NotEmpty(txRpc)
		s.Assert().Equal(txRpc.Hash().String(), mevTx.TxHash)
	}
}

func ForceDirToTestDirLocation() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	err := os.Chdir(dir)
	if err != nil {
		panic(err.Error())
	}
	return dir
}
