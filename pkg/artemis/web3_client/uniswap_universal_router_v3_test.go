package web3_client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func (s *Web3ClientTestSuite) TestUniversalRouterV3ExactIn() {
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
	hashStr := "0x3247555a5dbc877ade17c4b49362bc981af5fb5064e0b3cbd91411e085fe3093"
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

	for _, val := range subCmds.Commands {
		fmt.Println(val.Command)
	}
}

func (s *Web3ClientTestSuite) TestUniversalRouterV3ExactOut() {
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
	hashStr := "0xf99ac4237df313794747759550db919b37d7c8a67d4a7e12be8f5cbaacd51376"
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

	for _, val := range subCmds.Commands {
		fmt.Println(val.Command)
	}
}
