package transformations

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
	v1 "k8s.io/api/apps/v1"
)

type YamlWriterTestSuite struct {
	base.TestSuite
	y YamlFileIO
}

func (s *YamlWriterTestSuite) SetupTest() {
	s.y = YamlFileIO{}
}

func (s *YamlWriterTestSuite) TestDecodeK8sWorkloadDir() {
	fp := "deployment.yaml"
	jsonBytes, err := s.y.ReadYamlConfig(fp)
	s.Require().Nil(err)
	var d *v1.Deployment
	err = json.Unmarshal(jsonBytes, &d)
	s.Require().Nil(err)
	s.Assert().NotEmpty(d)
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./",
		FnIn:        "deployment.yaml",
		FnOut:       "write_deployment.yaml",
		Env:         "",
	}
	err = s.y.WriteYamlConfig(p, jsonBytes)
	s.Require().Nil(err)

}

func TestYamlWriterTestSuite(t *testing.T) {
	suite.Run(t, new(YamlWriterTestSuite))
}
