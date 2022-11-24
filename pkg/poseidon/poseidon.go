package poseidon

import (
	"strings"

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

	Protocol   string `json:"protocol"`
	Network    string `json:"network"`
	ClientType string `json:"clientType"`
	ClientName string `json:"clientName"`

	CompressionType string
}

func (b *BucketRequest) GetZstdFn() string {
	return b.GetCompressedBucketKey() + ".tar.zst"
}

func (b *BucketRequest) GetGzipFn() string {
	return b.GetCompressedBucketKey() + ".tar.gzip"
}

func (b *BucketRequest) GetBaseBucketKey() string {
	key := []string{strings.ToLower(b.Protocol), strings.ToLower(b.Network), strings.ToLower(b.ClientType), strings.ToLower(b.ClientName)}
	return strings.Join(key, ".")
}

func (b *BucketRequest) GetCompressedBucketKey() string {
	switch b.CompressionType {
	case "gzip":
		return b.GetGzipFn()
	default:
		return b.GetZstdFn()
	}
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
