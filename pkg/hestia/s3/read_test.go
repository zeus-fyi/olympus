package s3

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type S3TestSuite struct {
	suite.Suite
}

// TestRead, you'll need to set the secret values to run the test
func (t *S3TestSuite) TestRead() {
	ctx := context.Background()

	input := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String("test.txt"),
	}
	p := structs.Path{
		PackageName: "",
		DirIn:       "",
		DirOut:      "",
		Fn:          "local-text.txt",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	err := Read(ctx, p, input)
	t.Require().Nil(err)
}

func TestS3TestSuite(t *testing.T) {
	suite.Run(t, new(S3TestSuite))
}
