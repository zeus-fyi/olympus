package ethereum_slashing_protection_watermarking

import (
	"github.com/stretchr/testify/suite"
	hydra_base_test "github.com/zeus-fyi/olympus/hydra/api/test"
	"testing"
)

type WatermarkAttestationsWeb3SignerTestSuite struct {
	hydra_base_test.HydraBaseTestSuite
}

func (t *WatermarkAttestationsWeb3SignerTestSuite) TestAttestationWatermarker() {
	// TODO check source and target epochs
}

func (t *WatermarkAttestationsWeb3SignerTestSuite) TestFarFutureSigning() {
	// TODO check source and target epochs
}

func TestWatermarkAttestationsWeb3SignerTestSuite(t *testing.T) {
	suite.Run(t, new(WatermarkAttestationsWeb3SignerTestSuite))
}
