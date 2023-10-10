package apollo_status

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type StatusTestSuite struct {
	test_suites_base.TestSuite
}

func (t *StatusTestSuite) SetupTest() {
	t.InitLocalConfigs()
}

func (t *StatusTestSuite) TestStatusPolling() {

}

func TestStatusTestSuite(t *testing.T) {
	suite.Run(t, new(StatusTestSuite))
}
