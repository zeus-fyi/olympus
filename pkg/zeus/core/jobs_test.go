package zeus_core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/batch/v1"
)

type JobsTestSuite struct {
	K8TestSuite
}

var ctx = context.Background()

func (s *JobsTestSuite) TestGetJobsList() {
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "ephemeral"
	jl, err := s.K.GetJobsList(ctx, kns)
	s.Nil(err)
	s.Require().NotEmpty(jl)
}

func (s *JobsTestSuite) TestCreateJob() {
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "demo"
	// Delete
	// todo: create a job and add to sql
	job := &v1.Job{}
	j, err := s.K.CreateJob(ctx, kns, job)
	s.Nil(err)
	s.Require().NotEmpty(j)
}

func (s *JobsTestSuite) TestDeleteJob() {
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "demo"
	err := s.K.DeleteJob(ctx, kns, "jobName")
	s.Nil(err)
}

func TestJobsTestSuite(t *testing.T) {
	suite.Run(t, new(JobsTestSuite))
}
