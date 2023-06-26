package artemis_flashbots

import "github.com/metachris/flashbotsrpc"

func (s *FlashbotsTestSuite) TestGetFlashbotsBundleStats() {
	bundle := flashbotsrpc.FlashbotsGetBundleStatsParam{
		BlockNumber: "17565244",
		BundleHash:  "0xac5ca430bea0aeaf949308aaa2a7aa7ae424193174b97c292b0d1066ad685f3b",
	}
	resp, err := s.fb.GetBundleStats(ctx, bundle)
	s.Assert().Nil(err)
	s.Assert().NotNil(resp)
}
