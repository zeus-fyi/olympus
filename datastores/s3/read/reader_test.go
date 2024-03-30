package s3reader

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	s32 "github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_s3"
)

type S3ReadTestSuite struct {
	s32.S3TestSuite
}

// TestRead, you'll need to set the secret values to run the test
func (t *S3ReadTestSuite) TestRead() {
	ctx := context.Background()

	input := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String("text.txt"),
	}
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "",
		DirOut:      "",
		FnIn:        "local-text.txt",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	reader := NewS3ClientReader(t.S3)
	err := reader.Read(ctx, &p, input)
	t.Require().Nil(err)
}

// TestRead, you'll need to set the secret values to run the test
func (t *S3ReadTestSuite) TestReadOvh() {
	ctx := context.Background()

	input := &s3.GetObjectInput{
		Bucket: aws.String("zeusfyi"),
		Key:    aws.String("local-text.txt"),
	}
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "/Users/alex/go/Olympus/olympus/datastores/s3/read",
		DirOut:      "",
		FnIn:        "local-text.txt",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	reader := NewS3ClientReader(t.OvhS3)

	err := reader.Read(ctx, &p, input)
	t.Require().Nil(err)

}

func TestS3ReadTestSuite(t *testing.T) {
	suite.Run(t, new(S3ReadTestSuite))
}
