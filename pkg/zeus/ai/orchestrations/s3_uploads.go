package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	s3uploader "github.com/zeus-fyi/olympus/datastores/s3/upload"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/memfs"
)

func s3ws(ctx context.Context, cp *MbChildSubProcessParams, input *WorkflowStageIO) (*WorkflowStageIO, error) {
	if input == nil {
		log.Warn().Msg("s3ws: at least cp or input is nil or empty")
		return nil, fmt.Errorf("must have input to save s3 obj")
	}
	if err := errCheckStagedWfs(ctx, cp); err != nil {
		return nil, err
	}
	up := s3uploader.NewS3ClientUploader(athena.OvhS3Manager)
	p, err := stageNamePath(cp)
	if err != nil {
		log.Err(err).Msg("s3ws: failed to hash wsr io")
		return nil, err
	}
	b, err := json.Marshal(input)
	if err != nil {
		log.Err(err).Msg("s3ws: failed to upload wsr io")
		return nil, err
	}
	mfs := memfs.NewMemFs()
	err = mfs.MakeFileIn(p, b)
	if err != nil {
		log.Err(err).Msg("s3ws: failed to upload wsr io")
		return nil, err
	}
	kvs3 := &s3.PutObjectInput{
		Bucket: aws.String(FlowsBucketName),
		Key:    aws.String(p.FileOutPath()),
	}
	err = up.UploadFromInMemFsV2(ctx, p, kvs3, mfs)
	if err != nil {
		log.Err(err).Interface("p.FileOutPath()", p.FileOutPath()).Msg("s3ws: failed to upload wsr io")
		return nil, err
	}
	return input, err
}

func s3globalWf(ctx context.Context, cp *MbChildSubProcessParams, ue *artemis_entities.UserEntity) (*artemis_entities.UserEntity, error) {
	if err := errCheckGlobalWfs(ctx, cp, ue); err != nil {
		return nil, err
	}
	up := s3uploader.NewS3ClientUploader(athena.OvhS3Manager)
	p, err := globalWfEntityStageNamePath(cp, ue)
	if err != nil {
		log.Err(err).Msg("s3globalWf: failed to get filepath")
		return nil, err
	}
	b, err := json.Marshal(ue)
	if err != nil {
		log.Err(err).Interface("ue", ue).Msg("s3globalWf: failed to upload wsr io")
		return nil, err
	}
	mfs := memfs.NewMemFs()
	err = mfs.MakeFileIn(p, b)
	if err != nil {
		log.Err(err).Msg("s3globalWf: failed to upload wsr io")
		return nil, err
	}
	kvs3 := &s3.PutObjectInput{
		Bucket: aws.String(FlowsBucketName),
		Key:    aws.String(p.FileOutPath()),
	}
	err = up.UploadFromInMemFsV2(ctx, p, kvs3, mfs)
	if err != nil {
		log.Err(err).Interface("p.FileOutPath()", p.FileOutPath()).Msg("s3globalWf: failed to upload ue")
		return nil, err
	}
	return ue, err
}
