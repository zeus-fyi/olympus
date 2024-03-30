package test_suites_s3

import (
	"context"

	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/s3"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type S3TestSuite struct {
	test_suites_base.TestSuite
	S3    s3base.S3Client
	OvhS3 s3base.S3Client
}

func (s *S3TestSuite) SetupTest() {
	s.Tc = configs.InitLocalTestConfigs()
	s.SetupLocalDigitalOceanS3()
	s.SetupLocalOvhS3()
}

func (s *S3TestSuite) SetupLocalDigitalOceanS3() {
	ctx := context.Background()
	s3client, err := s3base.NewConnS3ClientWithStaticCreds(ctx, s.Tc.LocalS3SpacesKey, s.Tc.LocalS3SpacesSecret)
	s.Require().Nil(err)
	s.S3 = s3client
}

func (s *S3TestSuite) SetupLocalOvhS3() {
	ctx := context.Background()
	s3client, err := s3base.NewOvhConnS3ClientWithStaticCreds(ctx, s.Tc.OvhS3AccessKey, s.Tc.OvhS3SecretKey)
	s.Require().Nil(err)
	s.OvhS3 = s3client
}
