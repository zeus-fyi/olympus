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
	s3uploader "github.com/zeus-fyi/olympus/datastores/s3/upload"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/memfs"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

// s3ws uses workingRunCycleStagePath
func s3ws(ctx context.Context, cp *MbChildSubProcessParams, input *WorkflowStageIO) (*WorkflowStageIO, error) {
	if input == nil {
		log.Warn().Msg("s3ws: at least cp or input is nil or empty")
		return nil, fmt.Errorf("must have input to save s3 obj")
	}
	if err := errCheckStagedWfs(ctx, cp); err != nil {
		return nil, err
	}
	p, err := workingRunCycleStagePath(cp)
	if err != nil {
		log.Err(err).Msg("s3ws: failed to hash wsr io")
		return nil, err
	}
	b, err := json.Marshal(input)
	if err != nil {
		log.Err(err).Msg("s3ws: failed to upload wsr io")
		return nil, err
	}
	err = uploadFromInMemFs(ctx, b, p)
	if err != nil {
		log.Err(err).Interface("b", string(b)).Msg("failed to upload")
		return nil, err
	}
	return input, err
}

func s3wsCustomTaskName(ctx context.Context, cp *MbChildSubProcessParams, taskName string, input any) error {
	if input == nil {
		log.Warn().Msg("s3ws: at least cp or input is nil or empty")
		return fmt.Errorf("must have input to save s3 obj")
	}

	sn := cp.Tc.TaskName
	cp.Tc.TaskName = taskName

	log.Info().Str("taskName", taskName).Msg("s3wsCustomTaskName")
	if err := errCheckStagedWfs(ctx, cp); err != nil {
		return err
	}
	p, err := workingRunCycleStagePath(cp)
	if err != nil {
		log.Err(err).Msg("s3ws: failed to hash wsr io")
		return err
	}
	fmt.Println(p.FileOutPath(), "s3wsCustomTaskName")
	b, err := json.Marshal(input)
	if err != nil {
		log.Err(err).Msg("s3ws: failed to upload wsr io")
		return err
	}
	err = uploadFromInMemFs(ctx, b, p)
	if err != nil {
		log.Err(err).Interface("b", string(b)).Msg("failed to upload")
		return err
	}
	cp.Tc.TaskName = sn
	return nil
}

// s3globalWf uses real wf name "/%s/%s/%s", ogk, wfName, ue.Platform
func s3globalWf(ctx context.Context, cp *MbChildSubProcessParams, ue *artemis_entities.UserEntity) (*artemis_entities.UserEntity, error) {
	if err := errCheckGlobalWfs(ctx, cp, ue); err != nil {
		return nil, err
	}
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
	err = uploadFromInMemFs(ctx, b, p)
	if err != nil {
		log.Err(err).Interface("ue", ue).Msg("failed to upload")
		return nil, err
	}
	return ue, err
}

func S3GlobalOrgUpload(ctx context.Context, ou org_users.OrgUser, ue *artemis_entities.UserEntity) (*artemis_entities.UserEntity, error) {
	if ue == nil || ue.Platform == "" || len(ue.MdSlice) == 0 {
		return nil, fmt.Errorf("ue missing or field missing")
	}
	if ou.OrgID <= 0 {
		return nil, fmt.Errorf("org missing")
	}
	if err := s3SetupCheck(ctx); err != nil {
		return nil, err
	}
	p, err := globalOrgEntityStageNamePath(ou, ue, true)
	if err != nil {
		log.Err(err).Msg("s3globalWf: failed to get filepath")
		return nil, err
	}
	b, err := json.Marshal(ue)
	if err != nil {
		log.Err(err).Interface("ue", ue).Msg("s3globalWf: failed to upload wsr io")
		return nil, err
	}
	err = uploadFromInMemFs(ctx, b, p)
	if err != nil {
		log.Err(err).Interface("ue", ue).Msg("failed to upload")
		return nil, err
	}
	return ue, err
}

func S3WfRunImports(ctx context.Context, ou org_users.OrgUser, wfRunName string, ue *artemis_entities.UserEntity) (*artemis_entities.UserEntity, error) {
	if ue == nil || ue.Platform == "" || len(ue.MdSlice) == 0 {
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
	if len(ue.Nickname) <= 0 {
		if len(wfRunName) < 0 {
			return nil, fmt.Errorf("S3WfRunImports no nickname provided")
		}
		ue.Nickname = wfRunName
	}
	p := &filepaths.Path{
		DirIn:  fmt.Sprintf("/%s/%s/%s", ogk, wfRunName, ue.Platform),
		DirOut: fmt.Sprintf("/%s/%s/%s", ogk, wfRunName, ue.Platform),
		FnIn:   fmt.Sprintf("%s.json", ue.Nickname),
		FnOut:  fmt.Sprintf("%s.json", ue.Nickname),
	}
	log.Info().Interface("p.FileOutPath()", p.FileOutPath()).Msg("S3WfRunImports")
	b, err := json.Marshal(ue)
	if err != nil {
		log.Err(err).Interface("ue", ue).Msg("s3globalWf: failed to upload wsr io")
		return nil, err
	}
	err = uploadFromInMemFs(ctx, b, p)
	if err != nil {
		log.Err(err).Interface("ue", ue).Msg("failed to upload")
		return nil, err
	}
	return ue, err
}

func uploadFromInMemFs(ctx context.Context, b []byte, p *filepaths.Path) error {
	mfs := memfs.NewMemFs()
	err := mfs.MakeFileIn(p, b)
	if err != nil {
		log.Err(err).Msg("s3globalWf: failed to upload wsr io")
		return err
	}
	kvs3 := &s3.PutObjectInput{
		Bucket: aws.String(FlowsBucketName),
		Key:    aws.String(p.FileOutPath()),
	}
	up := s3uploader.NewS3ClientUploader(athena.OvhS3Manager)
	err = up.UploadFromInMemFsV2(ctx, p, kvs3, mfs)
	if err != nil {
		log.Err(err).Interface("p.FileOutPath()", p.FileOutPath()).Msg("s3globalWf: failed to upload ue")
		return err
	}
	return nil
}

func deleteFromS3(ctx context.Context, p *filepaths.Path) error {
	err := athena.OvhS3Manager.DeleteObject(ctx, FlowsBucketName, p.FileOutPath())
	if err != nil {
		log.Err(err).Interface("p.FileOutPath()", p.FileOutPath()).Msg("deleteFromInMemFs: failed to del")
		return err
	}
	return nil
}
