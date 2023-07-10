package artemis_trading_auxiliary

import (
	"fmt"

	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) TestCreateFbBundle() *AuxiliaryTradingUtils {
	ta := InitAuxiliaryTradingUtils(ctx, t.goerliNode, hestia_req_types.Goerli, t.acc)
	t.Require().NotEmpty(ta)
	toExchAmount := artemis_eth_units.GweiMultiple(100000)
	cmd := t.testEthToWETH(&ta, toExchAmount)
	// part 1 of bundle
	tx, err := ta.universalRouterCmdBuilder(ctx, cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	t.Require().Equal(toExchAmount, tx.Value())
	t.Require().Equal(1, len(ta.OrderedTxs))
	err = ta.CreateOrAddToFlashbotsBundle(cmd, "latest")
	t.Require().Nil(err)
	t.Require().NotEmpty(ta.Bundle.Txs)
	t.Require().Equal(1, len(ta.Bundle.Txs))
	t.Require().Equal(0, len(ta.OrderedTxs))

	// part 2 of bundle
	cmd = t.testExecV2Trade(&ta)
	tx, err = ta.universalRouterCmdBuilder(ctx, cmd)
	t.Require().NotEmpty(tx)
	t.Require().Equal(1, len(ta.OrderedTxs))
	err = ta.CreateOrAddToFlashbotsBundle(cmd, "latest")
	t.Require().Nil(err)
	t.Require().NotEmpty(ta.Bundle.Txs)
	t.Require().Equal(2, len(ta.Bundle.Txs))
	t.Require().Equal(0, len(ta.OrderedTxs))
	return &ta
}

func (t *ArtemisAuxillaryTestSuite) TestCallBundle() {
	ta := t.TestCreateFbBundle()
	t.Require().NotEmpty(ta)
	resp, err := ta.callFlashbotsBundle(ctx)
	t.Require().Nil(err)
	t.Require().NotNil(resp)

	fmt.Println("stateBlockNum", resp.StateBlockNumber)
	fmt.Println("bundleHash", resp.BundleHash)
	fmt.Println("gasFees", resp.GasFees)
	fmt.Println("totalGasUsed", resp.TotalGasUsed)

	t.Require().Equal(2, len(resp.Results))

	for _, sr := range resp.Results {
		t.Require().Emptyf(sr.Error, "error in result: %s", sr.Error)
	}
}
