package artemis_trading_auxiliary

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) testMockSandwichBundle() *AuxiliaryTradingUtils {
	ta := InitAuxiliaryTradingUtils(ctx, t.goerliNode, hestia_req_types.Goerli, t.acc)
	t.Require().NotEmpty(ta)
	toExchAmount := artemis_eth_units.GweiMultiple(100000)
	cmd := t.testEthToWETH(&ta, toExchAmount)
	// part 1 of bundle
	tx, err := ta.universalRouterCmdBuilder(ctx, cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	t.Require().Equal(toExchAmount, tx.Value())
	t.Require().Equal(1, len(ta.MevTxGroup.OrderedTxs))
	err = ta.CreateOrAddToFlashbotsBundle(cmd, "latest")
	t.Require().Nil(err)
	t.Require().NotEmpty(ta.Bundle.Txs)
	t.Require().Equal(1, len(ta.Bundle.Txs))
	t.Require().Equal(0, len(ta.MevTxGroup.OrderedTxs))

	// middle of bundle
	userTrader := InitAuxiliaryTradingUtils(ctx, t.goerliNode, hestia_req_types.Goerli, t.acc2)
	cmd = t.testEthToWETH(&userTrader, toExchAmount)
	tx, err = userTrader.universalRouterCmdBuilder(ctx, cmd)
	t.Require().NotEmpty(tx)
	err = ta.AddTxToBundleGroup(ctx, tx)
	t.Require().Nil(err)
	signer := types.LatestSignerForChainID(artemis_eth_units.NewBigInt(hestia_req_types.EthereumGoerliProtocolNetworkID))
	sender, err := signer.Sender(tx)
	t.Require().Nil(err)
	t.Require().Equal(t.acc2.Address().String(), sender.String())

	// part 3 of bundle
	cmd = t.testExecV2Trade(&ta)
	tx, err = ta.universalRouterCmdBuilder(ctx, cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	t.Require().Equal(2, len(ta.MevTxGroup.OrderedTxs))
	err = ta.CreateOrAddToFlashbotsBundle(cmd, "latest")
	t.Require().Nil(err)
	t.Require().NotEmpty(ta.Bundle.Txs)
	t.Require().Equal(3, len(ta.Bundle.Txs))
	t.Require().Equal(0, len(ta.MevTxGroup.OrderedTxs))
	return &ta
}
func (t *ArtemisAuxillaryTestSuite) TestSandwichCallBundle() {
	ta := t.testMockSandwichBundle()
	t.Require().NotEmpty(ta)
	resp, err := ta.callFlashbotsBundle(ctx)
	t.Require().Nil(err)
	t.Require().NotNil(resp)

	fmt.Println("stateBlockNum", resp.StateBlockNumber)
	fmt.Println("bundleHash", resp.BundleHash)
	fmt.Println("gasFees", resp.GasFees)
	fmt.Println("totalGasUsed", resp.TotalGasUsed)
	t.Require().Equal(3, len(resp.Results))

	for _, sr := range resp.Results {
		t.Require().Emptyf(sr.Error, "error in result: %s", sr.Error)
	}
}
