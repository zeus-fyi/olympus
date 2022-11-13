package zeus_ethereum

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	ares_test_suite "github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/ares"
)

type AresZeusEthereumTestSuite struct {
	ares_test_suite.AresTestSuite
}

var ctx = context.Background()

func (t *AresZeusEthereumTestSuite) SetupTest() {
	t.SetupLocalTest()
}

func TestAresZeusTestSuite(t *testing.T) {
	suite.Run(t, new(AresZeusEthereumTestSuite))
}
