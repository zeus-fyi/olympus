package web3_client

import (
	"fmt"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func (s *Web3ClientTestSuite) TestWrapETHFuncs() {
	node := "https://virulent-alien-cloud.quiknode.pro/fa84e631e9545d76b9e1b1c5db6607fedf3cb654"
	err := s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, node, 17461070)
	s.Require().Nil(err)
	userAddr := s.LocalHardhatMainnetUser.Address()
	b, err := s.LocalHardhatMainnetUser.GetBalance(ctx, userAddr.String(), nil)
	s.Require().Nil(err)
	fmt.Println("ethBalance", b.String())
	routerRecipient := accounts.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")
	wethParams := WrapETHParams{
		Recipient: routerRecipient,
		AmountMin: EtherMultiple(10),
	}
	payable := &web3_actions.SendEtherPayload{
		TransferArgs: web3_actions.TransferArgs{
			Amount:    EtherMultiple(10),
			ToAddress: wethParams.Recipient,
		},

		GasPriceLimits: web3_actions.GasPriceLimits{},
	}
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
		Payable:  payable,
	}
	encCmd, err := ur.EncodeCommands(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(encCmd)

	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
	tx, err := uni.ExecUniswapUniversalRouterCmd(ur)
	s.Assert().NoError(err)
	s.Assert().NotNil(tx)

	endTokenBalance, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, WETH9ContractAddress, userAddr.String())
	s.Require().Nil(err)
	s.Assert().Equal(EtherMultiple(10).String(), endTokenBalance.String())
	fmt.Println("endTokenBalance", endTokenBalance.String())

	approveTx, err := s.LocalHardhatMainnetUser.ERC20ApproveSpender(ctx, WETH9ContractAddress, UniswapUniversalRouterAddressOld, EtherMultiple(1000))
	s.Require().Nil(err)
	s.Require().NotNil(approveTx)
	unwrapWETHParams := UnwrapWETHParams{
		Recipient: routerRecipient,
		AmountMin: Ether,
	}
	transferTxParams := web3_actions.SendContractTxPayload{
		SmartContractAddr: WETH9ContractAddress,
		SendEtherPayload: web3_actions.SendEtherPayload{
			TransferArgs: web3_actions.TransferArgs{
				ToAddress: accounts.HexToAddress(UniswapUniversalRouterAddressOld),
			},
		},
		ContractABI: MustLoadERC20Abi(),
		Params:      []interface{}{accounts.HexToAddress(UniswapUniversalRouterAddressOld), Ether},
	}
	transferTx, err := s.LocalHardhatMainnetUser.TransferERC20Token(ctx, transferTxParams)
	s.Require().Nil(err)
	s.Require().NotNil(transferTx)

	ur = UniversalRouterExecCmd{
		Commands: []UniversalRouterExecSubCmd{
			{
				Command:       UnwrapWETH,
				CanRevert:     false,
				Inputs:        nil,
				DecodedInputs: unwrapWETHParams,
			},
		},
		Deadline: deadline,
		Payable:  nil,
	}
	encCmd, err = ur.EncodeCommands(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(encCmd)
	tx, err = uni.ExecUniswapUniversalRouterCmd(ur)
	s.Assert().NoError(err)
	s.Assert().NotNil(tx)
}

// TODO finish test case
func (s *Web3ClientTestSuite) TestExecV2TradeMethodUR() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	node := "https://virulent-alien-cloud.quiknode.pro/fa84e631e9545d76b9e1b1c5db6607fedf3cb654"
	err := s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, node, 17461070)
	s.Require().Nil(err)
	userAddr := s.LocalHardhatMainnetUser.Address()
	amountIn := EtherMultiple(2000)
	err = s.LocalHardhatMainnetUser.SetERC20BalanceBruteForce(ctx, LinkTokenAddr, userAddr.String(), amountIn)
	s.Require().Nil(err)
	addr1 := accounts.HexToAddress(LinkTokenAddr)
	addr2 := accounts.HexToAddress(WETH9ContractAddress)
	startTokenBalance, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, LinkTokenAddr, userAddr.String())
	s.Require().Nil(err)
	fmt.Println("startTokenBalance", startTokenBalance.String())
	b, err := s.LocalHardhatMainnetUser.GetBalance(ctx, userAddr.String(), nil)
	s.Require().Nil(err)
	fmt.Println("balance", b.String())
	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
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
	pair, err := uni.PairToPrices(ctx, []accounts.Address{accounts.HexToAddress(LinkTokenAddr), accounts.HexToAddress(WETH9ContractAddress)})
	s.Require().Nil(err)
	s.Require().NotEmpty(pair)
	amountOut, err := pair.GetQuoteUsingTokenAddr(WETH9ContractAddress, amountIn)
	s.Require().Nil(err)
	fmt.Println("amountOut", amountOut.String())
	to := TradeOutcome{
		AmountIn:      amountIn,
		AmountInAddr:  accounts.HexToAddress(LinkTokenAddr),
		AmountFees:    nil,
		AmountOut:     amountOut,
		AmountOutAddr: accounts.HexToAddress(WETH9ContractAddress),
	}
	_, err = s.LocalHardhatMainnetUser.ERC20ApproveSpender(ctx, to.AmountInAddr.String(), UniswapUniversalRouterAddressOld, to.AmountIn)
	s.Require().Nil(err)

	_, err = s.LocalHardhatMainnetUser.ERC20ApproveSpender(ctx, to.AmountInAddr.String(), WETH9ContractAddress, to.AmountIn)
	s.Require().Nil(err)

	startTokenBalance,
		err = s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, LinkTokenAddr, userAddr.String())
	s.Require().Nil(err)
	fmt.Println("startTokenBalance", startTokenBalance.String())

	v2ExactInTrade := V2SwapExactInParams{
		AmountIn:      amountIn,
		AmountOutMin:  amountOut,
		Path:          []accounts.Address{addr1, addr2},
		To:            accounts.HexToAddress(pair.PairContractAddr),
		PayerIsSender: true,
	}
	deadline, _ := new(big.Int).SetString("1461501637330902918203684832716283019655932542975", 10)
	var ur = UniversalRouterExecCmd{
		Commands: []UniversalRouterExecSubCmd{
			{
				Command:       V2SwapExactIn,
				CanRevert:     false,
				Inputs:        nil,
				DecodedInputs: v2ExactInTrade,
			},
		},
		Deadline: deadline,
	}
	encCmd, err := ur.EncodeCommands(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(encCmd)

	fmt.Println("startTokenBalance", startTokenBalance)
	tx, err := uni.ExecUniswapUniversalRouterCmd(ur)
	s.Assert().NoError(err)
	s.Assert().NotNil(tx)
	endTokenBalance, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, LinkTokenAddr, userAddr.String())
	s.Require().Nil(err)
	fmt.Println("endTokenBalance", endTokenBalance)
}
