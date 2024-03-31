package s3uploader

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/memfs"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func (s *S3ClientUploader) UploadFromInMemFsV2(ctx context.Context, p filepaths.Path, s3KeyValue *s3.PutObjectInput, inMemFs memfs.MemFS) error {
	log.Debug().Msg("UploadFromInMemFsV2")
	f, err := inMemFs.Open(p.FileOutPath())
	if err != nil {
		log.Err(err).Msg("S3ClientUploader: UploadFromInMemFsV2: p.OpenFileInPath()")
		return err
	}
	defer f.Close()

	s3KeyValue.Key = aws.String(p.FileOutPath())
	s3KeyValue.Body = f
	uploader := manager.NewUploader(s.AwsS3Client)
	_, err = uploader.Upload(ctx, s3KeyValue, func(u *manager.Uploader) {
	})
	if err != nil {
		log.Err(err).Msg("S3ClientUploader: UploadFromInMemFsV2: uploader.Upload(ctx, s3KeyValue)")
		return err
	}
	return err
}
