package transformations

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	v1 "k8s.io/api/apps/v1"
)

type TransformationTestSuite struct {
	test_suites_base.TestSuite
	y YamlFileIO
}

func (s *TransformationTestSuite) SetupTest() {
	s.y = YamlFileIO{}
}

func (s *TransformationTestSuite) TestDecodeK8sWorkloadDir() {
	p := filepaths.Path{
		DirIn: "./temp",
	}
	err := s.y.ReadK8sWorkloadDir(p)
	s.Require().Nil(err)
	s.Assert().NotEmpty(s.y.Deployment)
	s.Assert().NotEmpty(s.y.Service)
}

func (s *TransformationTestSuite) TestDecodeK8sWorkload() {
	fp := "deployment.yaml"
	err := s.y.DecodeK8sWorkload(fp)
	s.Require().Nil(err)
	s.Assert().NotEmpty(s.y.Deployment)
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
