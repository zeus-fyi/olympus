package web3_client

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

func (s *Web3ClientTestSuite) TestRawdawgExecUniversalRouterWETHSwap() {
	node := "https://virulent-alien-cloud.quiknode.pro/fa84e631e9545d76b9e1b1c5db6607fedf3cb654"
	err := s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, node, 17461070)
	s.Require().Nil(err)
	rawDawgPayload, bc, err := LoadLocalRawdawgAbiPayload()
	s.Require().Nil(err)
	rawDawgPayload.GasLimit = 2000000
	rawDawgPayload.Params = []interface{}{}

	tx, err := s.LocalHardhatMainnetUser.DeploySmartContract(ctx, bc, rawDawgPayload)
	s.Require().Nil(err)
	s.Assert().NotNil(tx)

	rx, err := s.LocalHardhatMainnetUser.WaitForReceipt(ctx, tx.Hash())
	s.Assert().Nil(err)
	s.Assert().NotNil(rx)
	rawdawgAddr := rx.ContractAddress.String()
	RawDawgAddr = rawdawgAddr

	userAddr := s.LocalHardhatMainnetUser.Address()
	b, err := s.LocalHardhatMainnetUser.GetBalance(ctx, userAddr.String(), nil)
	s.Require().Nil(err)
	fmt.Println("ethBalance", b.String())

	bal := hexutil.Big{}
	bigInt := bal.ToInt()
	bigInt.Set(EtherMultiple(100))
	bal = hexutil.Big(*bigInt)
	err = s.LocalHardhatMainnetUser.SetBalance(ctx, RawDawgAddr, bal)
	s.Require().Nil(err)
	b, err = s.LocalHardhatMainnetUser.GetBalance(ctx, RawDawgAddr, nil)
	s.Require().Nil(err)
	fmt.Println("rawdawgEthBalance", b.String())

	routerRecipient := accounts.HexToAddress(rawdawgAddr)
	wethParams := WrapETHParams{
		Recipient: routerRecipient,
		AmountMin: Ether,
	}
	//payable := &web3_actions.SendEtherPayload{
	//	TransferArgs: web3_actions.TransferArgs{
	//		Amount:    Ether,
	//		ToAddress: wethParams.Recipient,
	//	},
	//	GasPriceLimits: web3_actions.GasPriceLimits{},
	//}
	deadline, _ := new(big.Int).SetString("1461501637330902918203684832716283019655932542975", 10)
	ur := UniversalRouterExecCmd{
		Commands: []UniversalRouterExecSubCmd{
			{
				Command:       WrapETH,
				CanRevert:     false,
				Inputs:        nil,
				DecodedInputs: wethParams,
			},
		},
		Deadline: deadline,
		Payable:  nil,
	}
	encCmd, err := ur.EncodeCommands(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(encCmd)

	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
	tx, err = uni.ExecRawdawgUniversalRouterCmd(ur, nil)
	s.Assert().NoError(err)
	s.Assert().NotNil(tx)

	endTokenBalance, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, WETH9ContractAddress, rawdawgAddr)
	s.Require().Nil(err)
	s.Assert().Equal(Ether.String(), endTokenBalance.String())
}
