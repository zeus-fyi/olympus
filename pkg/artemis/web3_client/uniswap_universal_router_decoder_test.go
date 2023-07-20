package web3_client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func (s *Web3ClientTestSuite) TestUniswapUniversalRouterDecoding() {
	// needs to get txs sent to UR
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
	hashStr := "0xb3be368bf5fe4b8cafc80d60cd86152d06ce71c409280b71c9ed87b1088d042c"
	tx, _, err := s.MainnetWeb3User.GetTxByHash(ctx, common.HexToHash(hashStr))
	s.Require().Nil(err)
	s.Require().NotNil(tx)
	mn, args, err := DecodeTxArgData(ctx, tx, uni.MevSmartContractTxMapUniversalRouterNew)
	s.Require().Nil(err)
	s.Require().NotEmpty(mn)
	s.Require().NotEmpty(args)
	subCmds, err := NewDecodedUniversalRouterExecCmdFromMap(args, artemis_oly_contract_abis.MustLoadNewUniversalRouterAbi())
	s.Require().Nil(err)
	s.Require().NotEmpty(subCmds)
}
