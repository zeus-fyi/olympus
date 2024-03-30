package poseidon

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type Poseidon struct {
	compression.Compression
	s3base.S3Client
	filepaths.Path
}

type BucketRequest struct {
	BucketName string `json:"bucketName"`
	BucketKey  string `json:"bucketKey,omitempty"`

	Protocol        string `json:"protocol"`
	Network         string `json:"network"`
	ClientType      string `json:"clientType"`
	ClientName      string `json:"clientName"`
	CompressionType string `json:"compressionType,omitempty"`
}

type S3BucketRequest struct {
	BucketName      string `json:"bucketName"`
	BucketKey       string `json:"bucketKey,omitempty"`
	CompressionType string `json:"compressionType,omitempty"`
	//zeus_common_types.CloudCtxNs
}

func GetBinBuildBucket(appName string) BucketRequest {
	appName = strings.ToLower(appName)
	b := BucketRequest{}
	b.BucketName = appName
	b.BucketKey = appName
	b.CompressionType = "tar.lz4"
	return b
}

func GetBinBuildBucketKey(appName string) string {
	appName = strings.ToLower(appName)
	b := BucketRequest{}
	b.BucketName = appName
	b.BucketKey = appName
	b.CompressionType = "tar.lz4"
	key := []string{b.BucketKey, b.CompressionType}
	return strings.Join(key, ".")
}

func (b *BucketRequest) GetBucketKey() string {
	key := []string{strings.ToLower(b.Protocol), strings.ToLower(b.Network), strings.ToLower(b.ClientType), strings.ToLower(b.ClientName)}
	if len(b.CompressionType) == 0 {
		key = append(key, "tar.lz4")
	}
	return strings.Join(key, ".")
}

func (p *Poseidon) CheckIfBucketKeyExists(ctx context.Context, b BucketRequest) (bool, error) {
	_, err := p.S3Client.AwsS3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(b.BucketName),
		Key:    aws.String(b.GetBucketKey()),
	})
	if err != nil {
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func NewPoseidon(s3Client s3base.S3Client) Poseidon {
	return Poseidon{
		compression.NewCompression(),
		s3Client,
		filepaths.Path{
			PackageName: "",
			DirIn:       "/data",
			DirOut:      "/data",
			FnIn:        "",
			FnOut:       "",
			Env:         "",
			FilterFiles: string_utils.FilterOpts{},
		},
	}
}

func NewS3Poseidon(s3Client s3base.S3Client) Poseidon {
	return Poseidon{
		compression.NewCompression(),
		s3Client,
		filepaths.Path{
			PackageName: "",
			DirIn:       "./",
			DirOut:      "./",
			FnIn:        "",
			FnOut:       "",
			Env:         "",
			FilterFiles: string_utils.FilterOpts{},
		},
	}
}
