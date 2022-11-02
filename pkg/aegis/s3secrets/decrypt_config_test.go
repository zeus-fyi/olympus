package s3secrets

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type S3SecretsTestSuite struct {
	test_suites.S3TestSuite
}

// TestRead, you'll need to set the secret values to run the test
func (t *S3SecretsTestSuite) TestRead() {
	//ctx := context.Background()
	//
	//input := &s3.GetObjectInput{
	//	Bucket: aws.String("zeus-fyi"),
	//	Key:    aws.String("test.txt"),
	//}
	//p := structs.Path{
	//	PackageName: "",
	//	DirIn:       "",
	//	DirOut:      "",
	//	Fn:          "local-text.txt",
	//	Env:         "",
	//	FilterFiles: string_utils.FilterOpts{},
	//}
}

func TestS3SecretsTestSuite(t *testing.T) {
	suite.Run(t, new(S3SecretsTestSuite))
}
