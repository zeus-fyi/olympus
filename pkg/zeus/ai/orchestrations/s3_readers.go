package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
)

func gs3wfs(ctx context.Context, cp *MbChildSubProcessParams) (*WorkflowStageIO, error) {
	if err := errCheckStagedWfs(ctx, cp); err != nil {
		return nil, err
	}
	p, err := stageNamePath(cp)
	if err != nil {
		log.Err(err).Msg("gs3wfs: failed to hash wsr io")
		return nil, err
	}
	br := poseidon.S3BucketRequest{
		BucketName: FlowsBucketName,
		BucketKey:  p.FileOutPath(),
	}
	pos := poseidon.NewS3PoseidonLinux(athena.OvhS3Manager)
	buf, err := pos.S3DownloadReadBytes(ctx, br)
	if err != nil {
		log.Err(err).Interface("fp", p.FileOutPath()).Msg("gs3wfs: S3DownloadReadBytes error")
		return nil, err
	}
	input := &WorkflowStageIO{}
	err = json.Unmarshal(buf.Bytes(), &input)
	if err != nil {
		log.Err(err).Msg("gs3wfs: S3DownloadReadBytes error")
		return nil, err
	}
	return input, err
}

func gs3globalWf(ctx context.Context, cp *MbChildSubProcessParams, ue *artemis_entities.UserEntity) (*artemis_entities.UserEntity, error) {
	if err := errCheckGlobalWfs(ctx, cp, ue); err != nil {
		return nil, err
	}
	if err := s3SetupCheck(ctx); err != nil {
		return nil, err
	}
	p, err := globalWfEntityStageNamePath(cp, ue)
	if err != nil {
		log.Err(err).Msg("gs3globalWf: failed to get filepath")
		return nil, err
	}
	br := poseidon.S3BucketRequest{
		BucketName: FlowsBucketName,
		BucketKey:  p.FileOutPath(),
	}
	pos := poseidon.NewS3PoseidonLinux(athena.OvhS3Manager)
	buf, err := pos.S3DownloadReadBytes(ctx, br)
	if err != nil {
		log.Err(err).Msg("gs3globalWf: S3DownloadReadBytes error")
		return nil, err
	}
	err = json.Unmarshal(buf.Bytes(), &ue)
	if err != nil {
		log.Err(err).Msg("gs3globalWf: S3DownloadReadBytes error")
		return nil, err
	}
	return ue, err
}

func GetS3GlobalOrg(ctx context.Context, ou org_users.OrgUser, ue *artemis_entities.UserEntity) (*artemis_entities.UserEntity, error) {
	if ue == nil {
		return nil, fmt.Errorf("must have cp MbChildSubProcessParams to createe s3 obj key name")
	}
	if ou.OrgID <= 0 {
		return nil, fmt.Errorf("must have org id to save s3 obj")
	}
	if err := s3SetupCheck(ctx); err != nil {
		return nil, err
	}
	p, err := globalOrgEntityStageNamePath(ou, ue, false)
	if err != nil {
		log.Err(err).Msg("gs3globalWf: failed to get filepath")
		return nil, err
	}
	br := poseidon.S3BucketRequest{
		BucketName: FlowsBucketName,
		BucketKey:  p.FileOutPath(),
	}
	pos := poseidon.NewS3PoseidonLinux(athena.OvhS3Manager)
	buf, err := pos.S3DownloadReadBytes(ctx, br)
	if err != nil {
		log.Err(err).Msg("gs3globalWf: S3DownloadReadBytes error")
		return nil, err
	}
	err = json.Unmarshal(buf.Bytes(), &ue)
	if err != nil {
		log.Err(err).Msg("gs3globalWf: S3DownloadReadBytes error")
		return nil, err
	}
	return ue, err
}
