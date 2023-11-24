package s3reader

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

type S3ClientReader struct {
	s3base.S3Client
}

func NewS3ClientReader(baseClient s3base.S3Client) S3ClientReader {
	return S3ClientReader{
		baseClient,
	}
}

func (s *S3ClientReader) CheckIfKeyExists(ctx context.Context, s3KeyValue *s3.GetObjectInput) (bool, error) {
	_, err := s.GetHeadObject(ctx, s3KeyValue)
	if err != nil {
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		log.Err(err).Msg("S3ClientReader, GetHeadObject(ctx, s3KeyValue)")
		return false, err
	}
	return true, nil
}

func (s *S3ClientReader) Read(ctx context.Context, p *filepaths.Path, s3KeyValue *s3.GetObjectInput) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	downloader := manager.NewDownloader(s.AwsS3Client)
	newFile, err := os.Create(p.FileInPath())
	if err != nil {
		log.Err(err).Msgf("S3ClientReader, os.Create(p.FileInPath()), path: %s", p.FileInPath())
		return err
	}
	defer newFile.Close()
	_, err = downloader.Download(ctx, newFile, s3KeyValue)
	if err != nil {
		log.Err(err).Msg("S3ClientReader, downloader.Download(ctx, newFile, s3KeyValue)")
		return err
	}
	return err
}

func (s *S3ClientReader) ReadBytes(ctx context.Context, p *filepaths.Path, s3KeyValue *s3.GetObjectInput) *bytes.Buffer {
	if p == nil {
		panic(errors.New("need to include a path"))
	}

	log.Info().Msg("S3ClientReader: downloading bucket object")
	buf := &bytes.Buffer{}
	downloader := manager.NewDownloader(s.AwsS3Client)
	downloader.Concurrency = 1

	w := FakeWriterAt{w: buf}
	_, err := downloader.Download(ctx, w, s3KeyValue)
	if err != nil {
		log.Err(err).Interface("bucket", s3KeyValue).Msg("S3ClientReader: download failed shutting down server")
		misc.DelayedPanic(err)
	}
	return buf
}

func (s *S3ClientReader) ReadBytesNoPanic(ctx context.Context, p *filepaths.Path, s3KeyValue *s3.GetObjectInput) (*bytes.Buffer, error) {
	if p == nil {
		panic(errors.New("need to include a path"))
	}

	log.Info().Msg("Zeus: S3ClientReader, downloading bucket object")
	buf := &bytes.Buffer{}

	downloader := manager.NewDownloader(s.AwsS3Client)
	downloader.Concurrency = 1

	w := FakeWriterAt{w: buf}
	_, err := downloader.Download(ctx, w, s3KeyValue)

	if err != nil {
		log.Err(err).Msg("ReadBytesNoPanic")
		return buf, nil
	}

	return buf, err
}

type FakeWriterAt struct {
	w io.Writer
}

func (fw FakeWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	// ignore 'offset' because we forced sequential downloads
	return fw.w.Write(p)
}

func (s *S3ClientReader) GetHeadObject(ctx context.Context, s3KeyValue *s3.GetObjectInput) (*s3.HeadObjectOutput, error) {
	ho, err := s.S3Client.AwsS3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: s3KeyValue.Bucket,
		Key:    s3KeyValue.Key,
	})
	return ho, err
}
