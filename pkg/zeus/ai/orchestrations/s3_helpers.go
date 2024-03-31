package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	s3uploader "github.com/zeus-fyi/olympus/datastores/s3/upload"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
	"github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/memfs"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

// TODO integrate reader + swaps testing
func s3ws(ctx context.Context, input *WorkflowStageIO) (*WorkflowStageIO, error) {
	if input == nil {
		log.Warn().Msg("SaveEvalResponseOutput: at least one input is nil or empty")
		return nil, nil
	}
	if input.Org.OrgID <= 0 {
		return nil, fmt.Errorf("must have org id to save s3 obj")
	}
	if input.WorkflowOverrides.WorkflowRunName == "" {
		return nil, fmt.Errorf("must have run name to save s3 obj")
	}
	var err error
	if athena.OvhS3Manager.AwsS3Client == nil {
		// TODO verify secrets
		var ps *aws_secrets.OAuth2PlatformSecret
		ps, err = aws_secrets.GetMockingbirdPlatformSecrets(context.Background(), org_users.NewOrgUserWithID(FlowsOrgID, 0), "s3-ovh-us-west-or")
		if err != nil {
			log.Err(err).Msg("s3ws: failed to save workflow io")
			return nil, err
		}
		athena.OvhS3Manager, err = s3base.NewOvhConnS3ClientWithStaticCreds(ctx, ps.S3AccessKey, ps.S3SecretKey)
		if err != nil {
			log.Err(err).Msg("s3ws: failed to save workflow io")
			return nil, err
		}
	}
	up := s3uploader.NewS3ClientUploader(athena.OvhS3Manager)
	ogk, err := artemis_entities.HashParams(input.Org.OrgID, nil)
	if err != nil {
		log.Err(err).Msg("s3ws: failed to hash wsr io")
		return nil, err
	}
	p := filepaths.Path{
		DirIn:  fmt.Sprintf("/%s", ogk),
		DirOut: fmt.Sprintf("/%s", ogk),
		FnIn:   input.WorkflowExecParams.WorkflowOverrides.WorkflowRunName,
		FnOut:  input.WorkflowExecParams.WorkflowOverrides.WorkflowRunName,
	}
	b, err := json.Marshal(input)
	if err != nil {
		log.Err(err).Msg("s3ws: failed to upload wsr io")
		return nil, err
	}
	mfs := memfs.NewMemFs()
	err = mfs.MakeFileIn(&p, b)
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
		log.Err(err).Msg("s3ws: failed to upload wsr io")
		return nil, err
	}
	return input, err
}

// TODO integrate reader + swaps testing
func gs3wfs(ctx context.Context, input *WorkflowStageIO) (*WorkflowStageIO, error) {
	if input == nil {
		log.Warn().Msg("SaveEvalResponseOutput: at least one input is nil or empty")
		return nil, nil
	}
	var err error
	if athena.OvhS3Manager.AwsS3Client == nil {
		// TODO verify secrets
		var ps *aws_secrets.OAuth2PlatformSecret
		ps, err = aws_secrets.GetMockingbirdPlatformSecrets(context.Background(), org_users.NewOrgUserWithID(FlowsOrgID, 0), FlowsS3Ovh)
		if err != nil {
			log.Err(err).Msg("s3ws: failed to save workflow io")
			return nil, err
		}
		athena.OvhS3Manager, err = s3base.NewOvhConnS3ClientWithStaticCreds(ctx, ps.S3AccessKey, ps.S3SecretKey)
		if err != nil {
			log.Err(err).Msg("s3ws: failed to save workflow io")
			return nil, err
		}
	}

	var br poseidon.S3BucketRequest
	pos := poseidon.NewS3PoseidonLinux(athena.OvhS3Manager)
	err = pos.S3ZstdDownloadAndDec(ctx, br)
	if err != nil {
		log.Err(err).Msg("g3ws: S3ZstdDownloadAndDec error")
	}
	return nil, err
}
