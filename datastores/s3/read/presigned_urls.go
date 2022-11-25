package s3reader

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

func (s *S3ClientReader) GeneratePresignedURL(ctx context.Context, s3KeyValue *s3.GetObjectInput) (string, error) {
	preSignClient := s3.NewPresignClient(s.AwsS3Client)
	downloadURL, err := DownloadURL(ctx, preSignClient, s3KeyValue)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("S3ClientReader, DownloadURL")
		return downloadURL, err
	}
	return downloadURL, err
}

func DownloadURL(ctx context.Context, client *s3.PresignClient, s3KeyValue *s3.GetObjectInput) (string, error) {
	expiration := time.Now().Add(time.Hour * 12)
	getObjectArgs := s3.GetObjectInput{
		Bucket:          s3KeyValue.Bucket,
		ResponseExpires: &expiration,
		Key:             s3KeyValue.Key,
	}
	res, err := client.PresignGetObject(context.Background(), &getObjectArgs)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("S3ClientReader, DownloadURL")
		return "", err
	}
	return res.URL, nil
}
