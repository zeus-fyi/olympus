package kronos_helix

func (t *KronosWorkerTestSuite) TestAlert() {
	InitPagerDutyAlertClient(t.Tc.PagerDutyApiKey)
	t.Require().NotNil(PdAlertClient)
}
