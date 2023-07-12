package artemis_flashbots

import (
	"github.com/metachris/flashbotsrpc"
)

func (s *FlashbotsTestSuite) TestGetFlashbotsBundleStats() {
	bundle := flashbotsrpc.FlashbotsGetBundleStatsParam{
		BlockNumber: "0x10C063C", // 17565244
		BundleHash:  "0x9f93055488f7b9db678c14c1c5056c3ea01ef91e35c4f5e4cbeb6d8eb434f32d",
	}
	_, err := s.fb.FlashbotsGetBundleStatsV2(s.fb.getPrivateKey(), bundle)
	s.Assert().Nil(err)
}
