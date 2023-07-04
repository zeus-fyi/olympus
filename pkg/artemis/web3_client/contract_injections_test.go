package web3_client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

func (s *Web3ClientTestSuite) TestRawDawgInjection() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	err := s.LocalHardhatMainnetUser.HardhatResetNetworkToBlock(ctx, "https://virulent-alien-cloud.quiknode.pro/fa84e631e9545d76b9e1b1c5db6607fedf3cb654", 17408822)
	rawDawgPayload, bc := artemis_oly_contract_abis.MustLoadRawdawgContractDeployPayload()
	rawDawgPayload.GasLimit = 2000000
	rawDawgPayload.Params = []interface{}{}

	tx, err := s.LocalHardhatMainnetUser.DeploySmartContract(ctx, bc, rawDawgPayload)
	s.Require().Nil(err)
	s.Assert().NotNil(tx)

	rx, err := s.LocalHardhatMainnetUser.WaitForReceipt(ctx, tx.Hash())
	s.Assert().Nil(err)
	s.Assert().NotNil(rx)
	// TODO figure out why I can't set code override but somehow the deployment works fine wtf, is the ownable a problem?
	//err = s.LocalHardhatMainnetUser.SetCodeOverride(ctx, rx.ContractAddress.String(), artemis_oly_contract_abis.RawdawgByteCode)
	//s.Require().Nil(err)
	rawdawgAddr := rx.ContractAddress.String()
	//rawdawgAddr := rx.ContractAddress.String()
	daiAddr := "0x6B175474E89094C44Da98b954EedeAC495271d0F"
	userEth, err := s.LocalHardhatMainnetUser.GetBalance(ctx, s.LocalHardhatMainnetUser.PublicKey(), nil)
	s.Require().Nil(err)
	fmt.Println("userEth", userEth.String())

	bal := hexutil.Big{}
	bigInt := bal.ToInt()
	bigInt.Set(Ether)
	bal = hexutil.Big(*bigInt)
	err = s.LocalHardhatMainnetUser.SetBalance(ctx, rawdawgAddr, bal)
	s.Require().Nil(err)

	rawDawgStartingEth, err := s.LocalHardhatMainnetUser.GetBalance(ctx, rawdawgAddr, nil)
	s.Require().Nil(err)
	fmt.Println("rawDawgStartingEth", rawDawgStartingEth.String())

	abiInfo := artemis_oly_contract_abis.MustLoadRawdawgAbi()
	owner, err := s.LocalHardhatMainnetUser.GetOwner(ctx, abiInfo, rawdawgAddr)
	s.Require().Nil(err)
	fmt.Println(owner.String())

	err = s.LocalHardhatMainnetUser.SetERC20BalanceBruteForce(ctx, daiAddr, rawdawgAddr, TenThousandEther)
	s.Require().Nil(err)

	rawDawgDaiBal, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, daiAddr, rawdawgAddr)
	s.Require().Nil(err)
	s.Require().Equal(TenThousandEther, rawDawgDaiBal)
	fmt.Println("daiBalance", rawDawgDaiBal.String())
	// DAI-WETH pair contract
	daiWETHPairContractAddr := "0xA478c2975Ab1Ea89e8196811F51A7B7Ade33eB11"
	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
	pair, err := uni.V2PairToPrices(ctx, []accounts.Address{accounts.HexToAddress(daiAddr), accounts.HexToAddress(WETH9ContractAddress)})
	s.Require().Nil(err)
	s.Require().NotEmpty(pair)
	s.Require().Equal(pair.PairContractAddr, daiWETHPairContractAddr)
	amountIn := EtherMultiple(2000)

	amountOut, err := pair.GetQuoteUsingTokenAddr(daiAddr, amountIn)
	s.Require().Nil(err)
	fmt.Println("amountOut", amountOut.String())

	to := &artemis_trading_types.TradeOutcome{
		AmountInAddr:  accounts.HexToAddress(daiAddr),
		AmountIn:      amountIn,
		AmountOutAddr: accounts.HexToAddress(WETH9ContractAddress),
		AmountOut:     amountOut,
	}

	pathSlice := []string{to.AmountInAddr.String(), to.AmountOutAddr.String()}
	amountsOut, err := uni.GetAmountsOut(nil, to.AmountIn, pathSlice)
	s.Require().Nil(err)
	fmt.Println("amountsOut", amountsOut)
	for _, v := range amountsOut {
		fmt.Println(v.String())
	}

	tx, err = uni.ExecSmartContractTradingSwap(rawdawgAddr, pair, to)
	s.Require().Nil(err)
	s.Require().NotNil(tx)

	rawDawgEndingEth, err := s.LocalHardhatMainnetUser.GetBalance(ctx, rawdawgAddr, nil)
	s.Require().Nil(err)
	fmt.Println("rawDawgStartingEth", rawDawgStartingEth.String())
	fmt.Println("rawDawgEndingEth", rawDawgEndingEth.String())

	rawDawgDaiBal, err = s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, daiAddr, rawdawgAddr)
	s.Require().Nil(err)
	fmt.Println("rawDawgDaiBal", rawDawgDaiBal.String())

	rawDawgWETHbal, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, WETH9ContractAddress, rawdawgAddr)
	s.Require().Nil(err)
	fmt.Println("rawDawgWETHbal", rawDawgWETHbal.String())
	s.Require().Equal(amountOut.String(), rawDawgWETHbal.String())
}
