package s3writer

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type S3UploaderTestSuite struct {
	test_suites.S3TestSuite
}

// TestRead, you'll need to set the secret values to run the test
func (t *S3UploaderTestSuite) TestUpload() {
	ctx := context.Background()

	input := &s3.PutObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String("test.txt"),
	}
	p := structs.Path{
		PackageName: "",
		DirIn:       "",
		DirOut:      "",
		FnIn:        "local-text.txt",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	reader := NewS3ClientUploader(t.S3)
	err := reader.Upload(ctx, &p, input)
	t.Require().Nil(err)
}

func TestS3ReadTestSuite(t *testing.T) {
	suite.Run(t, new(S3UploaderTestSuite))
}
