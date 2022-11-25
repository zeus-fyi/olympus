package poseidon_pkg

import (
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
)

var PoseidonReader s3reader.S3ClientReader

func InitPoseidonReader(baseClient s3base.S3Client) {
	PoseidonReader = s3reader.NewS3ClientReader(baseClient)
}
