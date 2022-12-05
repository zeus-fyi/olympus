package apollo_client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	test_base "github.com/zeus-fyi/olympus/test"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

var ctx = context.Background()

type ApolloClientTestSuite struct {
	test_suites_base.TestSuite
	ApolloTestClient Apollo
}

func (t *ApolloClientTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()
	t.ApolloTestClient = NewDefaultApolloClient(tc.Bearer)
	// t.ApolloTestClient = NewLocalApolloClient(tc.Bearer)
	// points working dir to inside /test
	test_base.ForceDirToTestDirLocation()
}

func TestApolloClientTestSuite(t *testing.T) {
	suite.Run(t, new(ApolloClientTestSuite))
}
