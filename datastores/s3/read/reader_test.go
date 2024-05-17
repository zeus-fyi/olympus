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
	filepaths2 "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

type S3ReadTestSuite struct {
	s32.S3TestSuite
}

var ctx = context.Background()

func (t *S3ReadTestSuite) TestReadV2() {
	sr := NewS3ClientReader(t.OvhS3)
	p := &filepaths2.Path{
		PackageName: "",
		DirIn:       "/debug/runs",
		DirOut:      "/debug/runs",
		FnIn:        "csv-analysis-064b2553-602a.json",
		FnOut:       "",
		Env:         "",
	}
	s3KeyValue := &s3.GetObjectInput{
		Bucket: aws.String("flows"),
		Key:    aws.String(p.FileInPath()),
	}
	b, err := sr.ReadBytesNoPanicV2(ctx, p, s3KeyValue)
	t.Require().Nil(err)
	t.Require().NotEmpty(b)
}

// TestRead, you'll need to set the secret values to run the test
func (t *S3ReadTestSuite) TestRead() {
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
