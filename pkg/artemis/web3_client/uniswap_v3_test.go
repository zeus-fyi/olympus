package web3_client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (s *Web3ClientTestSuite) TestUniversalRouterV3() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ForceDirToTestDirLocation()
	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)

	hashStr := "0x2fee5123917e4f178b46d5e094cb58296c1c39b3f4589be191338b2ce98dca62"
	tx, _, err := s.MainnetWeb3User.GetTxByHash(ctx, common.HexToHash(hashStr))
	s.Require().Nil(err)
	s.Require().NotNil(tx)
	mn, args, err := DecodeTxArgData(ctx, tx, uni.MevSmartContractTxMapV3)
	s.Require().Nil(err)
	s.Require().NotEmpty(mn)
	s.Require().NotEmpty(args)

	ss := SwapExactInputSingleArgs{}
	err = ss.Decode(ctx, args)
	s.Require().Nil(err)

	//scInfo := MustLoadUniswapV3RouterAbi()
}
