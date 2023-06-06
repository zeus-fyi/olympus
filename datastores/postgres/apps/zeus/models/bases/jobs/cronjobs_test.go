package jobs

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type CronJobsTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *CronJobsTestSuite) TestCronJobs() {
	job := NewCronJob()
	filepath := s.TestDirectory + "/mocks/test/cronjob.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)
	err = json.Unmarshal(jsonBytes, &job.K8sCronJob)
	s.Require().Nil(err)
	s.Require().NotEmpty(job.K8sCronJob)
}

func TestCronJobsTestSuite(t *testing.T) {
	suite.Run(t, new(CronJobsTestSuite))
}
