package athena_fastsync

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type AthenaFastSyncTestSuite struct {
	test_suites_base.TestSuite
}

func (t *AthenaFastSyncTestSuite) SetupTest() {
	t.InitLocalConfigs()
}

func TestAthenaFastSyncTestSuite(t *testing.T) {
	suite.Run(t, new(AthenaFastSyncTestSuite))
}
