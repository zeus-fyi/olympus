package poseidon

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	s3uploader "github.com/zeus-fyi/olympus/datastores/s3/upload"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

func (p *Poseidon) ZstdCompressAndUpload(ctx context.Context, br BucketRequest) error {
	ctx = context.WithValue(ctx, "func", "ZstdCompressAndUpload")
	err := p.ZstCompressDir(&p.Path)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("ZstdCompressAndUpload: ZstCompressDir")
		return err
	}
	err = p.UploadSnapshot(ctx, br)
	return err
}

func (p *Poseidon) TarCompressAndUpload(ctx context.Context, br BucketRequest) error {
	ctx = context.WithValue(ctx, "func", "TarCompressAndUpload")
	err := p.TarCompress(&p.Path)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("TarCompressAndUpload: TarCompress")
		return err
	}
	err = p.UploadSnapshot(ctx, br)
	return err
}

func (p *Poseidon) Lz4CompressAndUpload(ctx context.Context, br BucketRequest) error {
	ctx = context.WithValue(ctx, "func", "Lz4CompressAndUpload")
	err := p.Lz4CompressDir(&p.Path)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Lz4CompressAndUpload: Lz4CompressDir")
		return err
	}
	err = p.UploadSnapshot(ctx, br)
	return err
}

func (p *Poseidon) GzipCompressAndUpload(ctx context.Context, br BucketRequest) error {
	ctx = context.WithValue(ctx, "func", "GzipCompressAndUpload")
	err := p.GzipCompressDir(&p.Path)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GzipCompressAndUpload: GzipCompressDir")
		return err
	}
	err = p.UploadSnapshot(ctx, br)
	return err
}

func (p *Poseidon) UploadSnapshot(ctx context.Context, br BucketRequest) error {
	ctx = context.WithValue(ctx, "func", "UploadSnapshot")
	uploader := s3uploader.NewS3ClientUploader(p.S3Client)
	input := &s3.PutObjectInput{
		Bucket: aws.String(br.BucketName),
		Key:    aws.String(br.GetBucketKey()),
	}
	err := uploader.Upload(ctx, p.Path, input)
	return err
}

func (p *Poseidon) SyncUpload(ctx context.Context, br BucketRequest) error {
	ctx = context.WithValue(ctx, "func", "SyncUpload")
	br.CompressionType = "none"
	spacesFolderLocation := fmt.Sprintf("spaces-sfo3:zeus-fyi-snapshots/%s", br.GetBucketKey())
	cmd := exec.Command("rclone", "sync", "data", spacesFolderLocation)
	err := cmd.Run()
	if err != nil {
		log.Fatal().Msg("Poseidon: SyncUpload failed")
		misc.DelayedPanic(err)
	}
	return err
}
