package s3uploader

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"github.com/zeus-fyi/olympus/sandbox/chains"
)

type S3UploaderTestSuite struct {
	test_suites.S3TestSuite
}

// TestRead, you'll need to set the secret values to run the test
func (t *S3UploaderTestSuite) TestUploadSimple() {
	ctx := context.Background()

	fn := "text.txt"
	bucketName := "zeus-fyi"
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("test.txt"),
	}
	p := structs.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "",
		FnIn:        fn,
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	uploader := NewS3ClientUploader(t.S3)
	err := uploader.Upload(ctx, &p, input)
	t.Require().Nil(err)
}

// TestRead, you'll need to set the secret values to run the test
func (t *S3UploaderTestSuite) TestUpload() {
	chains.ChangeToChainDataDir()
	ctx := context.Background()

	fn := "geth.tar.zst"
	bucketName := "zeus-fyi"
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("fn"),
	}
	p := structs.Path{
		PackageName: "",
		DirIn:       "./ethereum/geth_zstd_cmp",
		DirOut:      "",
		FnIn:        fn,
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	uploader := NewS3ClientUploader(t.S3)
	err := uploader.Upload(ctx, &p, input)
	t.Require().Nil(err)
}

func TestS3ReadTestSuite(t *testing.T) {
	suite.Run(t, new(S3UploaderTestSuite))
}
