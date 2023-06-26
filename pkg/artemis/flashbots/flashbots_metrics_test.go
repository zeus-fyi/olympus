package artemis_flashbots

import "github.com/metachris/flashbotsrpc"

func (s *FlashbotsTestSuite) TestGetFlashbotsBundleStats() {
	bundle := flashbotsrpc.FlashbotsGetBundleStatsParam{
		BlockNumber: "",
		BundleHash:  "",
	}
	resp, err := s.fb.GetBundleStats(ctx, bundle)
	s.Assert().Nil(err)
	s.Assert().NotNil(resp)
}
