package web3_client

import (
	"fmt"
	"os"
	"path"
	"runtime"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func (s *Web3ClientTestSuite) TestUniswapMempoolFilter() {
	ForceDirToTestDirLocation()
	uni := InitUniswapV2Client(ctx, s.MainnetWeb3User)
	uni.printOn = true
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
	uni.ProcessTxs()
	count := len(uni.SwapExactTokensForTokensParamsSlice)
	count += len(uni.SwapTokensForExactTokensParamsSlice)
	count += len(uni.SwapExactETHForTokensParamsSlice)
	count += len(uni.SwapTokensForExactETHParamsSlice)
	count += len(uni.SwapExactTokensForETHParamsSlice)
	count += len(uni.SwapETHForExactTokensParamsSlice)
	fmt.Println("Total trades found", count)
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
