package kronos_helix

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type KronosInstructionsTestSuite struct {
	test_suites_base.TestSuite
}

func (t *KronosInstructionsTestSuite) SetupTest() {
	t.InitLocalConfigs()
}

func (t *KronosInstructionsTestSuite) TestSendAlert() {

}

func TestKronosInstructionsTestSuite(t *testing.T) {
	suite.Run(t, new(KronosInstructionsTestSuite))
}
