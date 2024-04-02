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
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

const (
	csvSrcGlobalLabel      = "csv:global:source"
	csvSrcGlobalMergeLabel = "csv:global:merge"
)

func csvGlobalRetLabel() string {
	return fmt.Sprintf("%s:ret", csvSrcGlobalMergeLabel)
}

func csvGlobalMergeRetLabel(rn string) string {
	return fmt.Sprintf("%s:%s", csvGlobalRetLabel(), rn)
}

func gs3wfs(ctx context.Context, cp *MbChildSubProcessParams) (*WorkflowStageIO, error) {
	if err := errCheckStagedWfs(ctx, cp); err != nil {
		return nil, err
	}
	p, err := workingRunCycleStagePath(cp)
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

// GetGlobalEntitiesFromRef need to track all global; "csv:source:global" is only value specified for now
func GetGlobalEntitiesFromRef(ctx context.Context, ou org_users.OrgUser, refs []artemis_entities.EntitiesFilter) ([]artemis_entities.UserEntity, error) {
	var gens []artemis_entities.UserEntity
	for _, ev := range refs {
		ue := &artemis_entities.UserEntity{
			Nickname: ev.Nickname,
			Platform: ev.Platform,
		}
		ue, err := GetS3GlobalOrg(ctx, ou, ue)
		if err != nil {
			log.Err(err).Msg("GetGlobalEntitiesFromRef: failed to select workflow io")
			return nil, err
		}
		// todo replace with map search to support array of label inputs
		if !artemis_entities.SearchLabelsForMatch(csvSrcGlobalLabel, *ue) {
			continue
		}
		if len(ue.MdSlice) > 0 {
			gens = append(gens, *ue)
		}
	}
	return gens, nil
}

func S3WfRunExport(ctx context.Context, ou org_users.OrgUser, wfRunName string, ue *artemis_entities.UserEntity) (*artemis_entities.UserEntity, error) {
	if ue == nil || ue.Platform == "" || ue.Nickname == "" {
		return nil, fmt.Errorf("ue missing or field missing")
	}
	if ou.OrgID <= 0 {
		return nil, fmt.Errorf("org missing")
	}
	if err := s3SetupCheck(ctx); err != nil {
		return nil, err
	}
	ogk, err := artemis_entities.HashParams(ou.OrgID, nil)
	if err != nil {
		log.Err(err).Msg("workingRunCycleStagePath: failed to hash wsr io")
		return nil, err
	}
	p := &filepaths.Path{
		DirIn:  fmt.Sprintf("/%s/%s/%s", ogk, wfRunName, ue.Platform),
		DirOut: fmt.Sprintf("/%s/%s/%s", ogk, wfRunName, ue.Platform),
		FnIn:   fmt.Sprintf("%s.json", ue.Nickname),
		FnOut:  fmt.Sprintf("%s.json", ue.Nickname),
	}
	log.Info().Interface("p.FileOutPath()", p.FileOutPath()).Msg("S3WfRunImports")
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
