package s3uploader

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
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

func (s *S3ClientUploader) Upload(ctx context.Context, p structs.Path, s3KeyValue *s3.PutObjectInput) error {
	f, err := p.OpenFileInPath()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("S3ClientUploader: p.OpenFileInPath()")
		return err
	}
	defer f.Close()
	s3KeyValue.Body = f
	uploader := manager.NewUploader(s.AwsS3Client)
	_, err = uploader.Upload(ctx, s3KeyValue, func(u *manager.Uploader) {
		u.LeavePartsOnError = true // Don't delete the parts if the upload fails.
	})
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("S3ClientUploader: uploader.Upload(ctx, s3KeyValue)")
		return err
	}
	return err
}
