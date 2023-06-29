package beacon_api

func (s *BeaconAPITestSuite) TestWebsocket() {
	SubscribeToEvent(ctx, s.Tc.QuiknodeStreamWsNode)
}
