package s3uploader

import (
	"context"
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

type S3ClientUploader struct {
	s3base.S3Client
}

func NewS3ClientUploader(baseClient s3base.S3Client) S3ClientUploader {
	return S3ClientUploader{
		baseClient,
	}
}

func (s *S3ClientUploader) Upload(ctx context.Context, p filepaths.Path, s3KeyValue *s3.PutObjectInput) error {
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
		log.Err(err).Msg("S3ClientUploader: uploader.Upload(ctx, s3KeyValue)")
		return err
	}
	return err
}

func (s *S3ClientUploader) UploadFromInMemFs(ctx context.Context, p filepaths.Path, s3KeyValue *s3.PutObjectInput, inMemFs memfs.MemFS) error {
	log.Ctx(ctx).Debug().Msg("UploadFromInMemFs")
	f, err := inMemFs.Open(p.FileOutPath())
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("S3ClientUploader: UploadFromInMemFs: p.OpenFileInPath()")
		return err
	}
	defer f.Close()

	s3KeyValue.Key = aws.String(p.FnOut)
	s3KeyValue.Body = f
	s3KeyValue.Metadata = p.Metadata

	uploader := manager.NewUploader(s.AwsS3Client)
	_, err = uploader.Upload(ctx, s3KeyValue, func(u *manager.Uploader) {
	})
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("S3ClientUploader: UploadFromInMemFs: uploader.Upload(ctx, s3KeyValue)")
		return err
	}
	return err
}

func (s *S3ClientUploader) CheckIfKeyExists(ctx context.Context, s3KeyValue *s3.PutObjectInput) (bool, error) {
	_, err := s.S3Client.AwsS3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: s3KeyValue.Bucket,
		Key:    s3KeyValue.Key,
	})
	if err != nil {
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *S3ClientUploader) UploadBuildBinary(ctx context.Context, p filepaths.Path, s3KeyValue *s3.PutObjectInput) error {
	f, err := p.OpenFileInPath()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("S3ClientUploader:  UploadBuildBinary, p.OpenFileInPath()")
		return err
	}
	defer f.Close()
	s3KeyValue.Body = f
	uploader := manager.NewUploader(s.AwsS3Client)
	_, err = uploader.Upload(ctx, s3KeyValue, func(u *manager.Uploader) {
		u.LeavePartsOnError = true // Don't delete the parts if the upload fails.
	})
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("S3ClientUploader: UploadBuildBinary, uploader.Upload(ctx, s3KeyValue)")
		return err
	}
	return err
}
