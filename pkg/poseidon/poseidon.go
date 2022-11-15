package poseidon

import (
	"strings"

	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type Poseidon struct {
	compression.Compression
	s3base.S3Client
	structs.Path
}

type BucketRequest struct {
	BucketName string `json:"bucketName"`

	Protocol   string `json:"protocol"`
	Network    string `json:"network"`
	ClientType string `json:"clientType"`
	ClientName string `json:"clientName"`
}

func (b *BucketRequest) CreateBucketKey() string {
	key := []string{strings.ToLower(b.Protocol), strings.ToLower(b.Network), strings.ToLower(b.ClientName), strings.ToLower(b.ClientType)}
	return strings.Join(key, ".")
}

func NewPoseidon() Poseidon {
	return Poseidon{
		compression.NewCompression(),
		s3base.NewS3ClientBase(),
		structs.Path{
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
