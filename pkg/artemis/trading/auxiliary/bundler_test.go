package artemis_trading_auxiliary

func (t *ArtemisAuxillaryTestSuite) TestExecV2TradeBundle() {
	ta, cmd, _ := t.TestExecV2TradeCall()
	bundle := ta.CreateFlashbotsBundle(cmd, "latest")
	t.Require().NotEmpty(bundle)

	//_, err = ta.universalRouterExecuteTx(ctx, tx)
	//t.Require().Nil(err)
	//fmt.Println("tx", tx.Hash().String())
}
