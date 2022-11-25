package s3reader

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cavaliergopher/grab/v3"
)

func (t *S3ReadTestSuite) TestGeneratePresignedURL() {
	ctx := context.Background()
	input := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String("text.txt"),
	}
	reader := NewS3ClientReader(t.S3)
	url, err := reader.GeneratePresignedURL(ctx, input)
	t.Require().Nil(err)
	t.Assert().NotEmpty(url)

	client := grab.NewClient()
	req, err := grab.NewRequest(".", url)
	t.Require().Nil(err)

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)
	timer := time.NewTicker(500 * time.Millisecond)
	defer timer.Stop()
	select {
	case <-timer.C:
		fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
			resp.BytesComplete(),
			resp.Size(),
			100*resp.Progress())

	case <-resp.Done:
		// download is complete
		err = resp.Err()
		t.Require().Nil(err)
	}
}
