package s3base

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type S3ClientTestSuite struct {
	test_suites_base.TestSuite
}

func (t *S3ClientTestSuite) SetupTest() {
	t.InitLocalConfigs()
}

// TestRead, you'll need to set the secret values to run the test
func (t *S3ClientTestSuite) TestConnection() {
	s3client, err := NewOvhConnS3ClientWithStaticCreds(ctx, t.Tc.OvhS3AccessKey, t.Tc.OvhS3SecretKey)
	t.Require().Nil(err)
	t.Require().NotNil(s3client)
}

func TestS3ClientTestSuite(t *testing.T) {
	suite.Run(t, new(S3ClientTestSuite))
}
