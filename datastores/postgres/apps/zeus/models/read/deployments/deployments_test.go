package read_deployments

import (
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type DeploymentReaderTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *DeploymentReaderTestSuite) TestSelectQueryName() {
	//ctx := context.Background()
	//qp := test.CreateTestQueryNameParams()
	//
	//deploymentValues := Deployment{}
	//err := deploymentValues.SelectDeploymentResource(ctx, qp)
	//s.Require().Nil(err)
	//s.Require().NotEmpty(deploymentValues)
}

func TestDeploymentReaderTestSuite(t *testing.T) {
	suite.Run(t, new(DeploymentReaderTestSuite))
}
