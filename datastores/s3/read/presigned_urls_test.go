package s3reader

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (t *S3ReadTestSuite) TestGeneratePresignedURL() {
	ctx := context.Background()
	input := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String("test.txt"),
	}
	reader := NewS3ClientReader(t.S3)
	url, err := reader.GeneratePresignedURL(ctx, input)
	t.Require().Nil(err)
	t.Assert().NotEmpty(url)
}
