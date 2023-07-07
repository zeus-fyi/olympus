package artemis_flashbots

import "github.com/metachris/flashbotsrpc"

func (s *FlashbotsTestSuite) TestGetFlashbotsBundleSubmission() {
	bundle := flashbotsrpc.FlashbotsSendBundleRequest{
		Txs:          []string{"0x1", "0x2", "0x3"},
		BlockNumber:  "0x10C063C", // 17565244
		MinTimestamp: nil,
		MaxTimestamp: nil,
		RevertingTxs: nil,
	}
	_, err := s.fb.SendBundle(ctx, bundle)
	s.Assert().Nil(err)
}
