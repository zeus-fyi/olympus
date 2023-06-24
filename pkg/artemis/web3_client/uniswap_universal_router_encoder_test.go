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
				CanRevert:     true,
				Inputs:        nil,
				DecodedInputs: v2ExactInTrade,
			},
		},
	}
	encCmd, err := ur.EncodeCommands(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(encCmd)
	s.Require().NotNil(encCmd.Commands)
	subCmd := UniversalRouterExecSubCmd{}
	for i, byteVal := range encCmd.Commands {
		err = subCmd.DecodeCommand(byteVal, encCmd.Inputs[i])
		s.Require().NoError(err)
		s.Assert().Equal(true, subCmd.CanRevert)
		s.Assert().Equal(V2SwapExactIn, subCmd.Command)
		decodedInputs := subCmd.DecodedInputs.(V2SwapExactInParams)
		s.Assert().Equal(v2ExactInTrade.Path, decodedInputs.Path)
		s.Assert().Equal(v2ExactInTrade.AmountIn.String(), decodedInputs.AmountIn.String())
		s.Assert().Equal(v2ExactInTrade.AmountOutMin.String(), decodedInputs.AmountOutMin.String())
		s.Assert().Equal(v2ExactInTrade.To, decodedInputs.To)
		s.Assert().Equal(v2ExactInTrade.PayerIsSender, decodedInputs.PayerIsSender)
	}
}
