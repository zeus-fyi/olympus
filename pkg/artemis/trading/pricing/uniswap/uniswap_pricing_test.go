package artemis_uniswap_pricing

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
)

type UniswapPricingTestSuite struct {
	test_suites_encryption.EncryptionTestSuite
}

func (s *UniswapPricingTestSuite) SetupTest() {

}

func TestUniswapPricingTestSuite(t *testing.T) {
	suite.Run(t, new(UniswapPricingTestSuite))
}
