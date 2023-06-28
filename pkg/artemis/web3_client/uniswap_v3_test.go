package web3_client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (s *Web3ClientTestSuite) TestSwapExactInputSingleArgs() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ForceDirToTestDirLocation()
	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)

	hashStr := "0x2fee5123917e4f178b46d5e094cb58296c1c39b3f4589be191338b2ce98dca62"
	tx, _, err := s.MainnetWeb3User.GetTxByHash(ctx, common.HexToHash(hashStr))
	s.Require().Nil(err)
	s.Require().NotNil(tx)
	mn, args, err := DecodeTxArgData(ctx, tx, uni.MevSmartContractTxMapV3SwapRouterV2)
	s.Require().Nil(err)
	s.Require().NotEmpty(mn)
	s.Require().NotEmpty(args)

	ss := SwapExactInputSingleArgs{}
	err = ss.Decode(ctx, args)
	s.Require().Nil(err)

	hashStr = "0x91ec4779027dd501a8cd2d752f4c40768a7b37df7023d02ad6e4e21eacf172f3"
	tx, _, err = s.MainnetWeb3User.GetTxByHash(ctx, common.HexToHash(hashStr))
	s.Require().Nil(err)
}

func (s *Web3ClientTestSuite) TestSwapExactTokensForTokens() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ForceDirToTestDirLocation()
	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)

	hashStr := "0x91ec4779027dd501a8cd2d752f4c40768a7b37df7023d02ad6e4e21eacf172f3"
	tx, _, err := s.MainnetWeb3User.GetTxByHash(ctx, common.HexToHash(hashStr))
	s.Require().Nil(err)
	s.Require().NotNil(tx)
	mn, args, err := DecodeTxArgData(ctx, tx, uni.MevSmartContractTxMapV3SwapRouterV2)
	s.Require().Nil(err)
	s.Require().NotEmpty(mn)
	s.Require().NotEmpty(args)

	ss := SwapExactTokensForTokensParamsV3{}
	err = ss.Decode(ctx, args)
	s.Require().Nil(err)
	s.Assert().NotEmpty(ss)
}

func (s *Web3ClientTestSuite) TestDecodeMulticall() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ForceDirToTestDirLocation()
	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)

	hashStr := "0x5179d5430f7217356829c913036b187d910cd71080f6cb30c2d3ae96671a2bff"
	tx, _, err := s.MainnetWeb3User.GetTxByHash(ctx, common.HexToHash(hashStr))
	s.Require().Nil(err)
	s.Require().NotNil(tx)
	mn, args, err := DecodeTxArgData(ctx, tx, uni.MevSmartContractTxMapV3SwapRouterV2)
	s.Require().Nil(err)
	s.Require().NotEmpty(mn)
	s.Require().NotEmpty(args)
	mc := Multicall{}
	err = mc.Decode(ctx, args)
	s.Require().Nil(err)

	for _, data := range mc.Data {
		mn, args, err = DecodeTxData(ctx, data, uni.MevSmartContractTxMapV3SwapRouterV2.Abi, uni.MevSmartContractTxMapV3SwapRouterV2.Filter)
		s.Require().Nil(err)
		s.Require().NotEmpty(mn)
		s.Require().NotEmpty(args)
		in := ExactInputParams{}
		if mn == exactInput {
			err = in.Decode(ctx, args)
			s.Require().Nil(err)
			fmt.Println(in.TokenFeePath.TokenIn.String())
			fmt.Println(in.TokenFeePath.GetEndToken().String())
		}
	}
}
