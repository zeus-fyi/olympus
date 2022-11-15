package s3writer

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type S3ClientUploader struct {
	s3base.S3Client
}

func NewS3ClientUploader(baseClient s3base.S3Client) S3ClientUploader {
	return S3ClientUploader{
		baseClient,
	}
}

func (s *S3ClientUploader) Upload(ctx context.Context, p *structs.Path, s3KeyValue *s3.PutObjectInput) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	uploader := manager.NewUploader(s.AwsS3Client)
	newFile, err := os.Create(p.FnIn)
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = uploader.Upload(ctx, s3KeyValue)
	if err != nil {
		return err
	}
	return err
}
