package beacon_api

import "time"

func (s *BeaconAPITestSuite) TestWebsocket() {
	timestampChan := make(chan time.Time)
	TriggerWorkflowOnNewBlockHeaderEvent(ctx, s.Tc.QuikNodeStreamWsNode, timestampChan)
}
