package twillio

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type TwillioTestSuite struct {
	test_suites_base.TestSuite
}

func (t *TwillioTestSuite) SetupTest() {
	t.InitLocalConfigs()
}

func (t *TwillioTestSuite) TestStatusPolling() {

}

func TestTwillioTestSuite(t *testing.T) {
	suite.Run(t, new(TwillioTestSuite))
}
