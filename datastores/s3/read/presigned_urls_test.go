package s3reader

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cavaliergopher/grab/v3"
)

// TODO write all keys:presigned url & then pack it
func (t *S3ReadTestSuite) TestReadKeys() {
	ctx := context.Background()

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String("zeus-fyi-snapshots"),
		Prefix: aws.String("ethereum.mainnet.exec.client.standard.geth/"),
	}
	reader := NewS3ClientReader(t.S3)
	r, err := reader.ReadKeys(ctx, input)
	t.Require().Nil(err)
	t.Assert().NotEmpty(r)

	input.ContinuationToken = r.NextContinuationToken

	tmp := r.Contents[999]
	r2, err := reader.ReadKeys(ctx, input)
	t.Require().Nil(err)
	t.Assert().NotEmpty(r2)
	t.Assert().NotEqual(tmp, r2.Contents[0])
}

func (t *S3ReadTestSuite) TestGeneratePresignedURL() {
	ctx := context.Background()
	input := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi-test"),
		Key:    aws.String("test"),
	}
	reader := NewS3ClientReader(t.S3)
	url, err := reader.GeneratePresignedURL(ctx, input)
	t.Require().Nil(err)
	t.Assert().NotEmpty(url)
	fmt.Println(url)
}

func (t *S3ReadTestSuite) TestDownloadFromPresignedUrl() {
	ctx := context.Background()
	input := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi-test"),
		Key:    aws.String("gfdamnit.json"),
	}
	reader := NewS3ClientReader(t.S3)
	preSignedUrl, err := reader.GeneratePresignedURL(ctx, input)
	t.Require().Nil(err)
	t.Assert().NotEmpty(preSignedUrl)
	client := grab.NewClient()
	req, err := grab.NewRequest(".", preSignedUrl)
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
