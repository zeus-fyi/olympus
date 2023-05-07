package apollo_client

import (
	"context"
	"encoding/json"
	"io"
	"os"
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

const (
	UniSwapV2 = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
)

func (t *ApolloClientTestSuite) TestMempoolPrettyPrint() {
	jsonFile, err := os.Open("mempool/response.json")
	t.Require().NoError(err)
	defer jsonFile.Close()

	// Read the JSON file into a byte slice
	byteValue, err := io.ReadAll(jsonFile)
	t.Require().NoError(err)

	// Create an empty interface to hold the decoded JSON data
	var data interface{}

	// Unmarshal the JSON data into the interface
	err = json.Unmarshal(byteValue, &data)
	t.Require().NoError(err)

	// Marshal the interface back into pretty-printed JSON
	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	t.Require().NoError(err)

	// Write the new pretty-printed JSON back to the file
	err = os.WriteFile("mempool/response.json", prettyJSON, 0644)
	t.Require().NoError(err)
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
