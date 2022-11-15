package poseidon

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
)

func (p *Poseidon) Download(ctx context.Context, br BucketRequest) error {
	ctx = context.WithValue(ctx, "func", "GzipDownloadAndDec")

	downloader := s3reader.NewS3ClientReader(p.S3Client)
	input := &s3.GetObjectInput{
		Bucket: aws.String(br.BucketName),
		Key:    aws.String(br.CreateBucketKey()),
	}
	err := downloader.Read(ctx, &p.Path, input)
	if err != nil {
		return err
	}
	return err
}

func (p *Poseidon) ZstdDownloadAndDec(ctx context.Context, br BucketRequest) error {
	ctx = context.WithValue(ctx, "func", "ZstdDownloadAndDec")
	err := p.Download(ctx, br)
	if err != nil {
		return err
	}
	err = p.ZstdDecompress(&p.Path)
	return err
}

func (p *Poseidon) GzipDownloadAndDec(ctx context.Context, br BucketRequest) error {
	ctx = context.WithValue(ctx, "func", "GzipDownloadAndDec")
	err := p.Download(ctx, br)
	if err != nil {
		return err
	}
	err = p.GzipDecompress(&p.Path)
	return err
}
