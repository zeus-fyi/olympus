package s3reader

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
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

func (s *S3ClientReader) Read(ctx context.Context, p *structs.Path, s3KeyValue *s3.GetObjectInput) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	downloader := manager.NewDownloader(s.AwsS3Client)
	newFile, err := os.Create(p.FileInPath())
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = downloader.Download(ctx, newFile, s3KeyValue)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("S3ClientReader, downloader.Download(ctx, newFile, s3KeyValue)")
		return err
	}
	return err
}

func (s *S3ClientReader) ReadBytes(ctx context.Context, p *structs.Path, s3KeyValue *s3.GetObjectInput) *bytes.Buffer {
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
		log.Fatal().Msg("Zeus: S3ClientReader, download failed shutting down server")
		misc.DelayedPanic(err)
	}
	return buf
}

type FakeWriterAt struct {
	w io.Writer
}

func (fw FakeWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	// ignore 'offset' because we forced sequential downloads
	return fw.w.Write(p)
}
