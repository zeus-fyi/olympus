package web3_client

import (
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
)

func (s *Web3ClientTestSuite) TestUniversalRouterEncodeCommandByte() {
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
	ur := UniversalRouterExecCmd{
		Commands: []UniversalRouterExecSubCmd{
			{
				Command:       V2SwapExactIn,
				CanRevert:     false,
				Inputs:        nil,
				DecodedInputs: v2ExactInTrade,
			},
		},
	}
	encCmd, err := ur.EncodeCommands(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(encCmd)
	var cmdByte uint8
	subCmd := UniversalRouterExecSubCmd{}
	for _, byteVal := range encCmd.Commands {
		_, cmdByte, err = subCmd.DecodeCmdByte(byteVal)
		s.Require().NoError(err)
		s.Assert().Equal(uint8(V2_SWAP_EXACT_IN), cmdByte)
	}
}
