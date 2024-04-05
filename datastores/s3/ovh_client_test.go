package s3base

import (
	"context"
	"fmt"
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

	//pf := "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434aa7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"
	o, err := s3client.ListAllItemsInBucket(ctx, "flows")
	t.Require().Nil(err)
	t.Require().NotEmpty(o)
	//
	fmt.Println(len(o), "len")
	//for _, v := range o {
	//	err = s3client.DeleteObject(ctx, "flows", v)
	//	t.Require().Nil(err)
	//
	//}
}

func TestS3ClientTestSuite(t *testing.T) {
	suite.Run(t, new(S3ClientTestSuite))
}
