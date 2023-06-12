package web3_client

import (
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func (s *Web3ClientTestSuite) TestExecV2TradeMethodUR() {
	addr1 := accounts.HexToAddress(LidoSEthAddr)
	addr2 := accounts.HexToAddress(WETH9ContractAddress)
	v2ExactInTrade := V2SwapExactInParams{
		AmountIn:      big.NewInt(1000000000000000000),
		AmountOutMin:  big.NewInt(0),
		Path:          []accounts.Address{addr1, addr2},
		To:            accounts.HexToAddress(UniswapUniversalRouterAddress),
		PayerIsSender: true,
	}
	// convert to command
	var ur = UniversalRouterExecCmd{
		Commands: []UniversalRouterExecSubCmd{
			{
				Command:       V2SwapExactIn,
				CanRevert:     true,
				Inputs:        nil,
				DecodedInputs: v2ExactInTrade,
			},
		},
		Deadline: new(big.Int).SetUint64(1000000000000000000),
	}
	encCmd, err := ur.EncodeCommands(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(encCmd)

	uni := InitUniswapClient(ctx, s.HostedHardhatMainnetUser)
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
	err = s.HostedHardhatMainnetUser.HardHatResetNetwork(ctx, s.Tc.HardhatNode, 0)
	s.Require().Nil(err)
	tx, err := uni.ExecUniswapUniversalRouterCmd(ur)
	s.Require().NoError(err)
	s.Require().NotNil(tx)
}
