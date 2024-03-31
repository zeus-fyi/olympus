package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
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
