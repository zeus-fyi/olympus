package s3secrets

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func (s *S3Secrets) PullS3AndDecryptAndUnGzipToInMemFs(ctx context.Context, p *filepaths.Path, unzipDir string, bucketObj *s3.GetObjectInput) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	reader := s3reader.NewS3ClientReader(s.S3Client)
	err := reader.Read(ctx, p, bucketObj)
	if err != nil {
		return err
	}
	err = s.DecryptAndUnGzipToInMemFs(p, unzipDir)
	return err
}
