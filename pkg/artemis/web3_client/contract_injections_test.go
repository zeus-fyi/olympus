package web3_client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (s *Web3ClientTestSuite) TestRawDawgInjection() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	s.LocalHardhatMainnetUser.MustInjectRawDawg()
	bal := hexutil.Big{}
	bigInt := bal.ToInt()
	bigInt.Set(Ether)
	bal = hexutil.Big(*bigInt)
	err := s.LocalHardhatMainnetUser.SetBalance(ctx, RawDawgAddr, bal)
	s.Require().Nil(err)

	rawDawgBal, err := s.LocalHardhatMainnetUser.GetBalance(ctx, RawDawgAddr, nil)
	s.Require().Nil(err)
	s.Require().Equal(Ether, rawDawgBal)

	abiInfo := MustLoadRawdawgAbi()
	owner, err := s.LocalHardhatMainnetUser.GetOwner(ctx, abiInfo, RawDawgAddr)
	s.Require().Nil(err)
	fmt.Println(owner.String())

	daiAddr := "0x6b175474e89094c44da98b954eedeac495271d0f"
	err = s.LocalHardhatMainnetUser.SetERC20BalanceBruteForce(ctx, daiAddr, RawDawgAddr, TenThousandEther)
	s.Require().Nil(err)
	// DAI-USDC pair contract
	daiUsdcPairContractAddr := "0xa478c2975ab1ea89e8196811f51a7b7ade33eb11"
	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
	pair, err := uni.GetPairContractPrices(ctx, daiUsdcPairContractAddr)
	s.Require().Nil(err)
	s.Require().NotEmpty(pair)

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

}
