package s3

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func Read(ctx context.Context, p structs.Path, s3KeyValue *s3.GetObjectInput) error {
	awsS3Client, err := ConnectS3Session(ctx)
	if err != nil {
		return err
	}
	downloader := manager.NewDownloader(awsS3Client)
	newFile, err := os.Create(p.Fn)
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = downloader.Download(ctx, newFile, s3KeyValue)
	if err != nil {
		return err
	}
	return err
}
