package ethereum_slashing_protection_watermarking

import (
	"github.com/stretchr/testify/suite"
	hydra_base_test "github.com/zeus-fyi/olympus/hydra/api/test"
	"testing"
)

type WatermarkerWeb3SignerTestSuite struct {
	hydra_base_test.HydraBaseTestSuite
}

func (t *WatermarkerWeb3SignerTestSuite) TestBlockProposalWatermarker() {
	// TODO check slot
}

func (t *WatermarkerWeb3SignerTestSuite) TestAttestationWatermarker() {
	// TODO check source and target epochs
}

func TestWatermarkerWeb3SignerTestSuite(t *testing.T) {
	suite.Run(t, new(WatermarkerWeb3SignerTestSuite))
}
