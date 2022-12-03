package test_suites_s3

import (
	"context"

	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/s3"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type S3TestSuite struct {
	test_suites_base.TestSuite
	S3 s3base.S3Client
}

func (s *S3TestSuite) SetupTest() {
	s.SetupLocalDigitalOceanS3()
}

func (s *S3TestSuite) SetupLocalDigitalOceanS3() {
	s.Tc = configs.InitLocalTestConfigs()

	ctx := context.Background()
	s3client, err := s3base.NewConnS3ClientWithStaticCreds(ctx, s.Tc.LocalS3SpacesKey, s.Tc.LocalS3SpacesSecret)
	s.Require().Nil(err)
	s.S3 = s3client
}
