package transformations

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
	v1 "k8s.io/api/apps/v1"
)

type TransformationTestSuite struct {
	base.TestSuite
	y YamlReader
}

func (s *TransformationTestSuite) SetupTest() {
	s.y = YamlReader{}
}

func (s *TransformationTestSuite) TestDeploymentYamlParsing() {
	fp := "deployment.yaml"
	jsonBytes, err := s.y.ReadYamlConfig(fp)
	s.Require().Nil(err)
	var d *v1.Deployment
	err = json.Unmarshal(jsonBytes, &d)
	s.Require().Nil(err)
	s.Assert().NotEmpty(d)
}

func TestTransformationTestSuite(t *testing.T) {
	suite.Run(t, new(TransformationTestSuite))
}
