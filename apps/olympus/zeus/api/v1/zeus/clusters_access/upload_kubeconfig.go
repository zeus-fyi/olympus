package zeus_v1_clusters_api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	s3uploader "github.com/zeus-fyi/olympus/datastores/s3/upload"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type CreateOrUpdateKubeConfigsRequest struct {
}

func CreateOrUpdateKubeConfigsHandler(c echo.Context) error {
	request := new(CreateOrUpdateKubeConfigsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateKubeConfig(c)
}

var (
	AgeEnc  = encryption.Age{}
	KeysCfg = auth_startup.AuthConfig{}
)

func (t *CreateOrUpdateKubeConfigsRequest) CreateOrUpdateKubeConfig(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	if ou.OrgID == 0 {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	//fileResp, err := DecompressAndEncryptUserKubeConfigsWorkload(c)
	//if err != nil {
	//	log.Err(err).Msg("CreateOrUpdateKubeConfig: DecompressUserKubeConfigsWorkload")
	//	return err
	//}
	//err = EncAndUpload(c.Request().Context(), ou.OrgID, fileResp, AgeEnc)
	//if err != nil {
	//	log.Err(err).Msg("CreateOrUpdateKubeConfig: EncAndUpload")
	//	return c.JSON(http.StatusInternalServerError, nil)
	//}
	return c.JSON(http.StatusNotImplemented, nil)
}

func EncAndUpload(ctx context.Context, orgID int, in bytes.Buffer, ageEnc encryption.Age, kcf authorized_clusters.K8sClusterConfig) error {
	fn := fmt.Sprintf("%s.kube.tar.gz", kcf.GetHashedKey(orgID))
	bucketName := "zeus-fyi"
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fn),
	}
	p := filepaths.Path{
		DirIn:  "./",
		DirOut: "./",
		FnOut:  fn,
	}
	inMemFsEnc := memfs.NewMemFs()
	err := ageEnc.EncryptItem(inMemFsEnc, &p, in.Bytes())
	if err != nil {
		log.Err(err).Msg("CreateOrUpdateKubeConfig: EncryptItem")
		return err
	}
	uploader := s3uploader.NewS3ClientUploader(KeysCfg.GetS3BaseClient())
	err = uploader.UploadFromInMemFs(ctx, p, input, inMemFsEnc)
	if err != nil {
		log.Err(err).Msg("CreateOrUpdateKubeConfig: Upload")
		return err
	}
	inKey := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(p.FnOut),
	}
	inMemFsK8sConfig := auth_startup.ExtClustersRunDigitalOceanS3BucketObjAuthProcedure(ctx, KeysCfg, inKey)
	k := zeus_core.K8Util{}
	k.ConnectToK8sFromInMemFsCfgPath(inMemFsK8sConfig)
	ctxes, err := k.GetContexts()
	if err != nil {
		log.Err(err).Msg("CreateOrUpdateKubeConfig: GetContexts")
		return err
	}
	for name, _ := range ctxes {
		log.Info().Msgf("context name: %s", name)
	}
	cfgs := []authorized_clusters.K8sClusterConfig{kcf}
	err = authorized_clusters.InsertOrUpdateExtClusterConfigsUnique(ctx, org_users.NewOrgUserWithID(orgID, 0), cfgs)
	if err != nil {
		log.Err(err).Msg("CreateOrUpdateKubeConfig: InsertOrUpdateExtClusterConfigs")
		return err
	}
	return nil
}
