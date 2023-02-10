package ethereum_slashing_protection_watermarking

import (
	"github.com/stretchr/testify/suite"
	hydra_base_test "github.com/zeus-fyi/olympus/hydra/api/test"
	"testing"
)

type WatermarkBlockProposalsWeb3SignerTestSuite struct {
	hydra_base_test.HydraBaseTestSuite
}

func (t *WatermarkBlockProposalsWeb3SignerTestSuite) TestBlockProposalWatermarker() {
	// TODO check slot
}

func TestWatermarkBlockProposalsWeb3SignerTestSuite(t *testing.T) {
	suite.Run(t, new(WatermarkBlockProposalsWeb3SignerTestSuite))
}
