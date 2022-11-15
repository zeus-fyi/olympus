package poseidon

import (
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
