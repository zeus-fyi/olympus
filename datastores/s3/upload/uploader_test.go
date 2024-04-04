package s3uploader

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	s3test "github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_s3"
)

type S3UploaderTestSuite struct {
	s3test.S3TestSuite
}

// TestRead, you'll need to set the secret values to run the test
func (t *S3UploaderTestSuite) TestUploadZst() {
	ctx := context.Background()

	fn := "tmp.txt"
	bucketName := "zeusfyi"
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fn),
	}
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "/Users/alex/go/Olympus/olympus/datastores/s3/upload/",
		DirOut:      "./",
		FnIn:        fn,
	}

	uploader := NewS3ClientUploader(t.OvhS3)
	err := uploader.Upload(ctx, p, input)
	t.Require().Nil(err)
}

// TestRead, you'll need to set the secret values to run the test
func (t *S3UploaderTestSuite) TestUploadBuildBinary() {
	ctx := context.Background()

	fn := "artemis"
	bucketName := "zeus-fyi"
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fn),
	}
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "./ethereum/geth_zstd_cmp",
		DirOut:      "./",
		FnIn:        fn,
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	uploader := NewS3ClientUploader(t.S3)
	err := uploader.Upload(ctx, p, input)
	t.Require().Nil(err)
}

func TestS3ReadTestSuite(t *testing.T) {
	suite.Run(t, new(S3UploaderTestSuite))
}
