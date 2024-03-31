package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	"github.com/zeus-fyi/olympus/pkg/athena"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func s3SetupCheck(ctx context.Context) error {
	if athena.OvhS3Manager.AwsS3Client == nil {
		var ps *aws_secrets.OAuth2PlatformSecret
		ps, err := aws_secrets.GetMockingbirdPlatformSecrets(context.Background(), org_users.NewOrgUserWithID(FlowsOrgID, 0), FlowsS3Ovh)
		if err != nil {
			log.Err(err).Msg("s3ws: failed to save workflow io")
			return err
		}
		athena.OvhS3Manager, err = s3base.NewOvhConnS3ClientWithStaticCreds(ctx, ps.S3AccessKey, ps.S3SecretKey)
		if err != nil {
			log.Err(err).Msg("s3ws: failed to save workflow io")
			return err
		}
	}
	return nil
}

func stageNamePath(cp *MbChildSubProcessParams) (*filepaths.Path, error) {
	if cp == nil {
		return nil, fmt.Errorf("must have cp MbChildSubProcessParams to createe s3 obj key name")
	}
	if cp.Ou.OrgID <= 0 {
		return nil, fmt.Errorf("must have org id to save s3 obj")
	}
	if cp.WfExecParams.WorkflowOverrides.WorkflowRunName == "" {
		return nil, fmt.Errorf("must have run name to save s3 obj")
	}
	if len(cp.Tc.TaskName) <= 0 {
		return nil, fmt.Errorf("must have task name to save s3 obj")
	}
	// 1. wf-run-name
	// 2. wf-task-name
	// 3. wf-run-cycle
	// 4. wf-chunk
	wfRunName := cp.WfExecParams.WorkflowOverrides.WorkflowRunName
	ogk, err := artemis_entities.HashParams(cp.Ou.OrgID, nil)
	if err != nil {
		log.Err(err).Msg("s3ws: failed to hash wsr io")
		return nil, err
	}
	p := &filepaths.Path{
		DirIn:  fmt.Sprintf("/%s/%s/cycle/%d", ogk, wfRunName, cp.Wsr.RunCycle),
		DirOut: fmt.Sprintf("/%s/%s/cycle/%d", ogk, wfRunName, cp.Wsr.RunCycle),
		FnIn:   fmt.Sprintf("%s-%d.json", cp.Tc.TaskName, cp.Wsr.ChunkOffset),
		FnOut:  fmt.Sprintf("%s-%d.json", cp.Tc.TaskName, cp.Wsr.ChunkOffset),
	}
	return p, nil
}

// saves entities under the global wf name
func globalWfEntityStageNamePath(cp *MbChildSubProcessParams, ue *artemis_entities.UserEntity) (*filepaths.Path, error) {
	if cp == nil || ue == nil {
		return nil, fmt.Errorf("must have cp MbChildSubProcessParams to createe s3 obj key name")
	}
	if cp.Ou.OrgID <= 0 {
		return nil, fmt.Errorf("must have org id to save s3 obj")
	}
	if cp.WfExecParams.WorkflowTemplate.WorkflowName == "" {
		return nil, fmt.Errorf("must have wf name to save s3 obj entity")
	}
	if len(ue.Nickname) <= 0 {
		return nil, fmt.Errorf("entity have nickname to save s3 obj")
	}
	wfName := cp.WfExecParams.WorkflowTemplate.WorkflowName
	ogk, err := artemis_entities.HashParams(cp.Ou.OrgID, nil)
	if err != nil {
		log.Err(err).Msg("s3ws: failed to hash wsr io")
		return nil, err
	}
	p := &filepaths.Path{
		DirIn:  fmt.Sprintf("/%s/%s/%s", ogk, wfName, ue.Platform),
		DirOut: fmt.Sprintf("/%s/%s/%s", ogk, wfName, ue.Platform),
		FnIn:   fmt.Sprintf("%s.json", ue.Nickname),
		FnOut:  fmt.Sprintf("%s.json", ue.Nickname),
	}
	return p, nil
}

func errCheckStagedWfs(ctx context.Context, cp *MbChildSubProcessParams) error {
	if cp == nil {
		return fmt.Errorf("cp is nil")
	}
	if cp.Ou.OrgID <= 0 {
		log.Warn().Msg("gs3wfs: missing org id")
		return fmt.Errorf("gs3wfs: missing org id")
	}
	if len(cp.WfExecParams.WorkflowOverrides.WorkflowRunName) <= 0 {
		log.Warn().Msg("gs3wfs: missing run name")
		return fmt.Errorf("gs3wfs: missing run name")
	}
	return s3SetupCheck(ctx)
}

func errCheckGlobalWfs(ctx context.Context, cp *MbChildSubProcessParams, ue *artemis_entities.UserEntity) error {
	if cp == nil || ue == nil || len(ue.Nickname) <= 0 || len(ue.Platform) <= 0 {
		log.Warn().Msg("SaveEvalResponseOutput: least one of err: cp is nil or ue nickname/platform is empty")
		return fmt.Errorf("least one of err: cp is nil or ue nickname/platform is empty")
	}
	if cp.Ou.OrgID <= 0 {
		return fmt.Errorf("must have org id to save s3 obj")
	}
	if cp.WfExecParams.WorkflowOverrides.WorkflowRunName == "" {
		return fmt.Errorf("must have run name to save s3 obj")
	}
	return s3SetupCheck(ctx)
}
