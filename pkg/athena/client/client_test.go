package athena_client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
	test_base "github.com/zeus-fyi/olympus/test"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

var ctx = context.Background()

type AthenaClientTestSuite struct {
	base.TestSuite
	AthenaTestClient AthenaClient
}

func (t *AthenaClientTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()
	//t.ZeusTestClient = NewDefaultZeusClient(tc.Bearer)
	t.AthenaTestClient = NewLocalAthenaClient(tc.Bearer)
	// points working dir to inside /test
	test_base.ForceDirToTestDirLocation()
}

func TestAthenaClientTestSuite(t *testing.T) {
	suite.Run(t, new(AthenaClientTestSuite))
}
