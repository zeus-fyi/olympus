package ai_platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
)

// TODO integrate reader + swaps testing
func s3ws(ctx context.Context, input *WorkflowStageIO) (*WorkflowStageIO, error) {
	if input == nil {
		log.Warn().Msg("SaveEvalResponseOutput: at least one input is nil or empty")
		return nil, nil
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
	var br poseidon.S3BucketRequest
	pos := poseidon.NewS3PoseidonLinux(athena.OvhS3Manager)
	// TODO: replace with in memory bytes
	err = pos.S3ZstdCompressAndUpload(ctx, br)
	if err != nil {
		log.Err(err).Msg("sws: failed to SaveWorkflowIO")
		return nil, err
	}
	return nil, err
}

// ht, err := artemis_entities.HashWebRequestResultsAndParams(r.Ou, r.RouteInfo)
