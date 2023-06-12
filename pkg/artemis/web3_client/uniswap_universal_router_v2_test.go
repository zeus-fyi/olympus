package web3_client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func (s *Web3ClientTestSuite) TestUniversalRouterV2() {
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

	/*
		// 16606337
				trx_hash_03 = HexStr("0x889b34a27b730dd664cd71579b4310522c3b495fb34f17f08d1131c0cec651fa")
				expected_function_names_03 = ("WRAP_ETH", "V2_SWAP_EXACT_OUT", "UNWRAP_WETH")

			trx_hash_01 = HexStr("0x52e63b75f41a352ad9182f9e0f923c8557064c3b1047d1778c1ea5b11b979dd9")
			expected_function_names_01 = ("PERMIT2_PERMIT", "V2_SWAP_EXACT_IN")
	*/
	//
	// 0x889b34a27b730dd664cd71579b4310522c3b495fb34f17f08d1131c0cec651fa
	// 16591736
	hashStr := "0x889b34a27b730dd664cd71579b4310522c3b495fb34f17f08d1131c0cec651fa"
	tx, _, err := s.MainnetWeb3User.GetTxByHash(ctx, common.HexToHash(hashStr))
	s.Require().Nil(err)
	s.Require().NotNil(tx)
	mn, args, err := DecodeTxArgData(ctx, tx, uni.MevSmartContractTxMapUniversalRouter)
	s.Require().Nil(err)
	s.Require().NotEmpty(mn)
	s.Require().NotEmpty(args)
	subCmds, err := NewDecodedUniversalRouterExecCmdFromMap(args)
	s.Require().Nil(err)
	s.Require().NotEmpty(subCmds)

	node := "https://virulent-alien-cloud.quiknode.pro/fa84e631e9545d76b9e1b1c5db6607fedf3cb654"
	err = s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, node, 16591736)
	s.Require().Nil(err)
}
