package web3_client

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func (s *Web3ClientTestSuite) TestUniswapUniversalRouterDecoding() {
	// needs to get txs sent to UR
	// can lookup a tx hash on etherscan and get the input data
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
	tx := &types.Transaction{}
	mn, args, err := DecodeTxArgData(ctx, tx, uni.UniversalRouterAbi, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(mn)
	s.Require().NotEmpty(args)
}
