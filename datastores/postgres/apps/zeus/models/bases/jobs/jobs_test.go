package jobs

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type JobsTestSuite struct {
	test_suites_base.TestSuite
	TestDirectory string
}

func (s *JobsTestSuite) SetupTest() {
	s.TestDirectory = "."
}

func (s *JobsTestSuite) TestK8sToDBJobParsing() {
	jo := NewJob()
	filepath := path.Join(s.TestDirectory, "job.yaml")

	b, err := ReadYamlConfig(filepath)
	s.Require().Nil(err)
	s.Require().NotEmpty(b)

	err = json.Unmarshal(b, &jo.K8sJob)
	s.Require().Nil(err)

	err = jo.ConvertK8JobToDB()
	s.Require().Nil(err)

	s.Assert().NotEmpty(jo.Metadata)
	s.Require().Equal("example-job", jo.Metadata.Name.ChartSubcomponentValue)
}

func (s *JobsTestSuite) TestK8sJob() {
	job := NewJob()
	filepath := path.Join(s.TestDirectory, "job.yaml")
	jsonBytes, err := ReadYamlConfig(filepath)
	s.Require().Nil(err)
	err = json.Unmarshal(jsonBytes, &job.K8sJob)
	s.Require().Nil(err)
	s.Require().NotEmpty(job.K8sJob)
}

func ReadYamlConfig(filepath string) ([]byte, error) {
	// Open YAML file
	jsonByteArray, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	jsonBytes, err := yaml.YAMLToJSON(jsonByteArray)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return jsonBytes, err
	}
	return jsonBytes, err
}

func TestJobsTestSuite(t *testing.T) {
	suite.Run(t, new(JobsTestSuite))
}
