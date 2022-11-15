package athena_chain_snapshots

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BucketRequest struct {
	Network    string `json:"network"`
	ClientType string `json:"clientType"`
	ClientName string `json:"clientName"`
}

func (b *BucketRequest) CreateBucketKey() *s3.GetObjectInput {
	key := []string{strings.ToLower(b.Network), strings.ToLower(b.ClientType), strings.ToLower(b.ClientName)}
	bucketKey := strings.Join(key, ".")
	var input = &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String(bucketKey),
	}
	return input
}
