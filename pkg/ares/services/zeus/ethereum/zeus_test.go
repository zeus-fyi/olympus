package ethereum

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	ares_test_suite "github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/ares"
)

type AresZeusTestSuite struct {
	ares_test_suite.AresTestSuite
}

var ctx = context.Background()

func (t *AresZeusTestSuite) SetupTest() {
	t.SetupProdTest()
}

func TestAresZeusTestSuite(t *testing.T) {
	suite.Run(t, new(AresZeusTestSuite))
}
