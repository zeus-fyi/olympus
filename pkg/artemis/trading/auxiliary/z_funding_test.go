package artemis_trading_auxiliary

import (
	"context"
	"fmt"
	"sort"

	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_reporting "github.com/zeus-fyi/olympus/pkg/artemis/trading/reporting"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
)

// 	GetNextPermit2NonceFromContract(ctx, )

func (t *ArtemisAuxillaryTestSuite) TestMainnetGetNextPermit2NonceFromContract() {
	age := encryption.NewAge(t.Tc.LocalAgePkey, t.Tc.LocalAgePubkey)
	t.acc3 = initTradingAccount2(ctx, age)
	w3aMainnet := web3_client.NewWeb3Client(t.mainnetNode, &t.acc3)
	w3aMainnet.AddBearerToken(t.Tc.ProductionLocalTemporalBearerToken)
	//_, err := GetNextPermit2NonceFromContract(ctx, w3aMainnet, t.acc3.Address(), accounts.HexToAddress("0x285DB79fa7e0e89E822786F48A7c98C6c1dC1c7d"), artemis_trading_constants.UniswapUniversalRouterNewAddressAccount)
	//t.Require().Nil(err)
	//
	rw, err := artemis_reporting.GetRewardsHistory(ctx, artemis_reporting.RewardHistoryFilter{
		FromBlock:   17658962,
		TradeMethod: "any",
	})
	t.Require().Nil(err)
	total := artemis_eth_units.NewBigInt(0)
	totalWithoutNegatives := artemis_eth_units.NewBigInt(0)
	// Sort slice
	var historySlice []artemis_reporting.RewardsHistory
	for _, v := range rw.Map {
		total = artemis_eth_units.AddBigInt(total, v.ExpectedProfitAmountOut)
		if artemis_eth_units.IsXGreaterThanY(v.ExpectedProfitAmountOut, artemis_eth_units.NewBigInt(0)) {
			totalWithoutNegatives = artemis_eth_units.AddBigInt(totalWithoutNegatives, v.ExpectedProfitAmountOut)
		}
		rh := artemis_reporting.RewardsHistory{
			FailedCount:             v.FailedCount,
			AmountOutToken:          v.AmountOutToken,
			Count:                   v.Count,
			ExpectedProfitAmountOut: v.ExpectedProfitAmountOut,
		}
		historySlice = append(historySlice, rh)
	}

	sort.SliceStable(historySlice, func(i, j int) bool {
		// Descending order
		return historySlice[i].Count > historySlice[j].Count
	})

	negCount := 0
	var addresses []string
	for _, v := range historySlice {
		total1 := v.Count + v.FailedCount
		if v.Count < 2 || (float64(v.FailedCount) > float64(total1)*0.2) {
			// v.FailedCount is more than 10% of the total
			log.Info().Str("token", v.AmountOutToken.Name()).Str("address", v.AmountOutToken.Address.String()).Int("successCount", v.Count).Int("failureCount", v.FailedCount).Msg("failed count is more than 20% of the total")
			continue
		}

		if artemis_eth_units.IsXGreaterThanY(artemis_eth_units.NewBigInt(0), rw.Map[v.AmountOutToken.Address.String()].ExpectedProfitAmountOut) {
			log.Warn().Str("token", v.AmountOutToken.Name()).Str("address", v.AmountOutToken.Address.String()).Int("successCount", v.Count).Int("failureCount", v.FailedCount).Msg("expected profit is negative")
			negCount++
			continue
		}
		fmt.Println(
			"successCount", v.Count, "failedCount", v.FailedCount, "expProfits", v.ExpectedProfitAmountOut.String(),
			v.AmountOutToken.Name(), v.AmountOutToken.Address.String(),
			"num", v.AmountOutToken.TransferTax.Numerator.String(), "den", v.AmountOutToken.TransferTax.Denominator.String())
		addresses = append(addresses, v.AmountOutToken.Address.String())
	}

	if len(addresses) <= 65 {
		return
	}
	fmt.Println("total", len(addresses))
	offset := uint64(0)
	atMainnet := InitAuxiliaryTradingUtils(ctx, w3aMainnet)
	for _, addr := range addresses {
		allowance, aerr := w3aMainnet.ReadERC20Allowance(ctx, addr, t.acc3.Address().String(), artemis_trading_constants.Permit2SmartContractAddress)
		t.Require().Nil(aerr)
		fmt.Println("allowance", allowance.String())

		if artemis_eth_units.IsXLessThanEqZeroOrOne(allowance) {
			ctx = web3_actions.SetNonceOffset(context.Background(), offset)
			res, er := atMainnet.SetPermit2ApprovalForToken(ctx, addr)
			t.Require().Nil(er)
			t.Require().NotNil(res)
			offset += 1
		}
	}
	//
	//err = artemis_risk_analysis.SetTradingPermission(ctx, addresses, 1, true)
	//t.Assert().Nil(err)
}
func (t *ArtemisAuxillaryTestSuite) TestMainnetPermitAllowance() {
	age := encryption.NewAge(t.Tc.LocalAgePkey, t.Tc.LocalAgePubkey)
	t.acc3 = initTradingAccount2(ctx, age)
	w3aMainnet := web3_client.NewWeb3Client(t.mainnetNode, &t.acc3)
	atMainnet := InitAuxiliaryTradingUtils(ctx, w3aMainnet)

	addr := "0x285DB79fa7e0e89E822786F48A7c98C6c1dC1c7d"

	allowance, aerr := w3aMainnet.ReadERC20Allowance(ctx, "0x285DB79fa7e0e89E822786F48A7c98C6c1dC1c7d", t.acc3.Address().String(), artemis_trading_constants.Permit2SmartContractAddress)
	t.Require().Nil(aerr)
	fmt.Println("allowance", allowance.String())

	if artemis_eth_units.IsXLessThanEqZeroOrOne(allowance) {
		res, er := atMainnet.SetPermit2ApprovalForToken(ctx, addr)
		t.Require().Nil(er)
		t.Require().NotNil(res)
	}
}

func (t *ArtemisAuxillaryTestSuite) TestMainnetBal() {
	age := encryption.NewAge(t.Tc.LocalAgePkey, t.Tc.LocalAgePubkey)
	t.acc3 = initTradingAccount2(ctx, age)
	w3aMainnet := web3_client.NewWeb3Client(t.mainnetNode, &t.acc3)
	w3aMainnet.AddBearerToken(t.Tc.ProductionLocalTemporalBearerToken)
	atMainnet := InitAuxiliaryTradingUtils(ctx, w3aMainnet)
	token := getChainSpecificWETH(*atMainnet.w3c()).String()
	fmt.Println("token", token)
	t.Require().NotEmpty(w3aMainnet.Headers)

	t.Require().Equal("Bearer "+t.Tc.ProductionLocalTemporalBearerToken, w3aMainnet.Headers["Authorization"])
	t.Require().Equal(t.mainnetNode, atMainnet.nodeURL())
	t.Require().Equal(token, artemis_trading_constants.WETH9ContractAddress)

	bal, err := checkEthBalance(ctx, *atMainnet.w3c())
	t.Require().Nil(err)
	t.Require().NotNil(bal)

	fmt.Println("bal", bal.String())

	bal, err = CheckAuxWETHBalance(ctx, *atMainnet.w3c())
	t.Require().Nil(err)
	t.Require().NotNil(bal)
	fmt.Println("weth bal", bal.String())
}

func (t *ArtemisAuxillaryTestSuite) TestSetPermit2Mainnet() {
	age := encryption.NewAge(t.Tc.LocalAgePkey, t.Tc.LocalAgePubkey)
	t.acc3 = initTradingAccount2(ctx, age)
	w3aMainnet := web3_client.NewWeb3Client(t.mainnetNode, &t.acc3)
	w3aMainnet.AddBearerToken(t.Tc.ProductionLocalTemporalBearerToken)
	atMainnet := InitAuxiliaryTradingUtils(ctx, w3aMainnet)
	token := getChainSpecificWETH(*atMainnet.w3c()).String()
	fmt.Println("token", token)
	t.Require().NotEmpty(w3aMainnet.Headers)

	t.Require().Equal("Bearer "+t.Tc.ProductionLocalTemporalBearerToken, w3aMainnet.Headers["Authorization"])
	t.Require().Equal(t.mainnetNode, atMainnet.nodeURL())
	t.Require().Equal(token, artemis_trading_constants.WETH9ContractAddress)

	//approveTx, err := atMainnet.SetPermit2ApprovalForToken(ctx, token)
	//t.Require().Nil(err)
	//t.Require().NotEmpty(approveTx)
	//fmt.Println("approveTx", approveTx.Hash().String())
}

func (t *ArtemisAuxillaryTestSuite) TestFundAccount() {
	age := encryption.NewAge(t.Tc.LocalAgePkey, t.Tc.LocalAgePubkey)
	t.acc3 = initTradingAccount2(ctx, age)
	w3aMainnet := web3_client.NewWeb3Client("https://eth.zeus.fyi", &t.acc3)
	w3aMainnet.AddBearerToken(t.Tc.ProductionLocalTemporalBearerToken)
	atMainnet := InitAuxiliaryTradingUtils(ctx, w3aMainnet)
	token := getChainSpecificWETH(*atMainnet.w3c()).String()
	fmt.Println("token", token)
	t.Require().NotEmpty(w3aMainnet.Headers)

	t.Require().Equal("Bearer "+t.Tc.ProductionLocalTemporalBearerToken, w3aMainnet.Headers["Authorization"])
	t.Require().Equal(t.mainnetNode, atMainnet.nodeURL())
	t.Require().Equal(token, artemis_trading_constants.WETH9ContractAddress)

	bal, err := checkEthBalance(ctx, *atMainnet.w3c())
	t.Require().Nil(err)
	t.Require().NotNil(bal)

	fmt.Println("bal", bal.String())

	bal, err = CheckAuxWETHBalance(ctx, *atMainnet.w3c())
	t.Require().Nil(err)
	t.Require().NotNil(bal)
	fmt.Println("weth bal", bal.String())

	// 0.45 eth
	toExchAmount := artemis_eth_units.GweiMultiple(450000000)

	cmd := t.testEthToWETH(&atMainnet, toExchAmount)
	found := false
	for i, sc := range cmd.Commands {
		if i == 0 && sc.Command == artemis_trading_constants.WrapETH {
			found = true
			t.Require().NotNil(cmd.Payable.Amount)
			t.Require().Equal(artemis_trading_constants.UniswapUniversalRouterNewAddressAccount.String(), cmd.Payable.ToAddress.String())
			t.Require().Equal(toExchAmount.String(), cmd.Payable.Amount.String())
			t.Require().Equal(toExchAmount.String(), sc.DecodedInputs.(web3_client.WrapETHParams).AmountMin.String())
			t.Require().Equal(artemis_trading_constants.UniversalRouterSender, sc.DecodedInputs.(web3_client.WrapETHParams).Recipient.String())
		}
	}
	t.Require().True(found)
	ok, err := checkEthBalanceGreaterThan(ctx, *atMainnet.w3c(), toExchAmount)
	t.Require().Nil(err)
	t.Require().True(ok)

	tx, _, err := universalRouterCmdToTxBuilder(ctx, *atMainnet.w3c(), cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	// 122
	//executedTx, err := atMainnet.universalRouterExecuteTx(ctx, tx)
	//t.Require().Nil(err)
	//fmt.Println("executedTx", executedTx.Hash().String())
}
