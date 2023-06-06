package jobs

import (
	"encoding/json"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type CronJobsTestSuite struct {
	test_suites_base.TestSuite
	TestDirectory string
}

func (s *CronJobsTestSuite) SetupTest() {
	s.TestDirectory = "."
}

func (s *CronJobsTestSuite) TestCronJobs() {
	job := NewCronJob()
	filepath := path.Join(s.TestDirectory, "cronjob.yaml")
	jsonBytes, err := ReadYamlConfig(filepath)
	s.Require().Nil(err)
	err = json.Unmarshal(jsonBytes, &job.K8sCronJob)
	s.Require().Nil(err)
	s.Require().NotEmpty(job.K8sCronJob)
}

func TestCronJobsTestSuite(t *testing.T) {
	suite.Run(t, new(CronJobsTestSuite))
}
