package artemis_trading_auxiliary

func (t *ArtemisAuxillaryTestSuite) TestExecV2TradeBundle() {
	ta, cmd, _ := t.TestExecV2TradeCall()
	err := ta.CreateFlashbotsBundle(cmd, "latest")
	t.Require().Nil(err)

	//_, err = ta.universalRouterExecuteTx(ctx, tx)
	//t.Require().Nil(err)
	//fmt.Println("tx", tx.Hash().String())
}
