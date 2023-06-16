package zeus_core

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	v1 "k8s.io/api/batch/v1"
)

type CronJobsTestSuite struct {
	K8TestSuite
}

func (s *CronJobsTestSuite) TestGetJobsList() {
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "ephemeral"
	jl, err := s.K.GetCronJobsList(ctx, kns)
	s.Nil(err)
	s.Require().NotEmpty(jl)
}

func (s *CronJobsTestSuite) TestCreateJob() {
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "demo"
	// Delete
	// todo: create a job and add to sql
	job := &v1.CronJob{}
	j, err := s.K.CreateCronJob(ctx, kns, job)
	s.Nil(err)
	s.Require().NotEmpty(j)
}

func (s *CronJobsTestSuite) TestDeleteJob() {
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "demo"
	err := s.K.DeleteCronJob(ctx, kns, "jobName")
	s.Nil(err)
}

func TestCronJobsTestSuite(t *testing.T) {
	suite.Run(t, new(CronJobsTestSuite))
}
