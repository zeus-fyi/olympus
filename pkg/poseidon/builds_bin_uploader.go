package poseidon

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	s3uploader "github.com/zeus-fyi/olympus/datastores/s3/upload"
)

// UploadsBinCompressBuild you need to set the path in Poseidon, then it does everything else
func (p *Poseidon) UploadsBinCompressBuild(ctx context.Context, appName string) error {
	ctx = context.WithValue(ctx, "func", "UploadsBinCompressBuild")
	uploader := s3uploader.NewS3ClientUploader(p.S3Client)
	// Upload the binary
	keyValue := &s3.PutObjectInput{
		Bucket: aws.String("zeus.fyi"),
		Key:    aws.String(GetBinBuildBucketKey(appName)),
	}
	err := uploader.Upload(ctx, p.Path, keyValue)
	if err != nil {
		log.Err(err).Msg("Poseidon: UploadsBinCompressBuild")
		return err
	}
	return nil
}
