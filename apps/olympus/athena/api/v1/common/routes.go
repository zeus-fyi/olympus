package v1_common_routes

import (
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

var CommonManager ClientManager

type ClientManager struct {
	S3Client s3reader.S3ClientReader
	DataDir  filepaths.Path
}
