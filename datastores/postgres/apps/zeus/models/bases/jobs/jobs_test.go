package jobs

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type JobsTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *JobsTestSuite) TestK8sSecretYamlReaderAndK8sToDBCte() {
	job := NewJob()
	filepath := s.TestDirectory + "/mocks/test/job.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)
	err = json.Unmarshal(jsonBytes, &job.K8sJob)
	s.Require().Nil(err)
	s.Require().NotEmpty(job.K8sJob)
}

func TestJobsTestSuite(t *testing.T) {
	suite.Run(t, new(JobsTestSuite))
}
