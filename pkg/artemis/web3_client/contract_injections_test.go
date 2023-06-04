package web3_client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (s *Web3ClientTestSuite) TestRawDawgInjection() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	err := s.LocalHardhatMainnetUser.HardhatResetNetworkToBlock(ctx, "https://virulent-alien-cloud.quiknode.pro/fa84e631e9545d76b9e1b1c5db6607fedf3cb654", 17408822)
	s.Require().Nil(err)
	daiAddr := "0x6B175474E89094C44Da98b954EedeAC495271d0F"
	s.LocalHardhatMainnetUser.MustInjectRawDawg()
	userEth, err := s.LocalHardhatMainnetUser.GetBalance(ctx, s.LocalHardhatMainnetUser.PublicKey(), nil)
	s.Require().Nil(err)
	fmt.Println("userEth", userEth.String())

	rawDawgStartBal, err := s.LocalHardhatMainnetUser.GetBalance(ctx, RawDawgAddr, nil)
	s.Require().Nil(err)
	fmt.Println("rawDawgStartBal", rawDawgStartBal.String())
	bal := hexutil.Big{}
	bigInt := bal.ToInt()
	bigInt.Set(Ether)
	bal = hexutil.Big(*bigInt)
	err = s.LocalHardhatMainnetUser.SetBalance(ctx, RawDawgAddr, bal)
	s.Require().Nil(err)

	rawDawgBal, err := s.LocalHardhatMainnetUser.GetBalance(ctx, RawDawgAddr, nil)
	s.Require().Nil(err)
	s.Require().Equal(Ether, rawDawgBal)

	abiInfo := MustLoadRawdawgAbi()
	owner, err := s.LocalHardhatMainnetUser.GetOwner(ctx, abiInfo, RawDawgAddr)
	s.Require().Nil(err)
	fmt.Println(owner.String())

	err = s.LocalHardhatMainnetUser.SetERC20BalanceBruteForce(ctx, daiAddr, RawDawgAddr, TenThousandEther)
	s.Require().Nil(err)

	rawDawgDaiBal, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, daiAddr, RawDawgAddr)
	s.Require().Nil(err)
	s.Require().Equal(TenThousandEther, rawDawgDaiBal)
	fmt.Println("daiBalance", rawDawgDaiBal.String())
	// DAI-USDC pair contract
	daiUsdcPairContractAddr := "0xA478c2975Ab1Ea89e8196811F51A7B7Ade33eB11"
	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
	pair, err := uni.PairToPrices(ctx, []accounts.Address{accounts.HexToAddress(daiAddr), accounts.HexToAddress(WETH9ContractAddress)})
	s.Require().Nil(err)
	s.Require().NotEmpty(pair)
	s.Require().Equal(pair.PairContractAddr, daiUsdcPairContractAddr)
	amountIn := EtherMultiple(2000)
	amountOut, err := pair.GetQuoteUsingTokenAddr(daiAddr, amountIn)
	s.Require().Nil(err)
	fmt.Println("amountOut", amountOut.String())

	to := &TradeOutcome{
		AmountInAddr:  accounts.HexToAddress(daiAddr),
		AmountIn:      amountIn,
		AmountOutAddr: accounts.HexToAddress(WETH9ContractAddress),
		AmountOut:     amountOut,
	}
	_, err = uni.ExecSmartContractTradingSwap(pair, to)
	s.Require().Nil(err)

	rawDawgDaiBal, err = s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, daiAddr, RawDawgAddr)
	s.Require().Nil(err)
	fmt.Println("rawDawgDaiBal", rawDawgDaiBal.String())

	rawDawgWETHbal, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, WETH9ContractAddress, RawDawgAddr)
	s.Require().Nil(err)
	fmt.Println("rawDawgWETHbal", rawDawgWETHbal.String())
	s.Require().Equal(amountOut.String(), rawDawgWETHbal.String())
}
