package s3secrets

import (
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

type S3Secrets struct {
	compression.Compression
	encryption.Age
	s3reader.S3ClientReader
	memfs.MemFS
}

func NewS3Secrets(c compression.Compression, a encryption.Age, s3r s3reader.S3ClientReader) S3Secrets {
	s3secrets := S3Secrets{
		Compression:    c,
		Age:            a,
		S3ClientReader: s3r,
		MemFS:          memfs.NewMemFs(),
	}
	return s3secrets
}
