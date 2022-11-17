package zeus_demo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	ares_test_suite "github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/ares"
)

type AresDemoTestSuite struct {
	ares_test_suite.AresTestSuite
}

var ctx = context.Background()

func (t *AresDemoTestSuite) SetupTest() {

	t.SetupDemoProdUserTest()

	//t.SetupLocalTest()
}

func TestAresZeusTestSuite(t *testing.T) {
	suite.Run(t, new(AresDemoTestSuite))
}
