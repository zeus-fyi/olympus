package poseidon

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	s3uploader "github.com/zeus-fyi/olympus/datastores/s3/upload"
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
		Key:    aws.String(br.GetCompressedBucketKey()),
	}
	p.FnIn = br.GetBaseBucketKey()
	err := uploader.Upload(ctx, p.Path, input)
	return err
}
