package zeus_v1_clusters_api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/ext_clusters"
	s3uploader "github.com/zeus-fyi/olympus/datastores/s3/upload"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/athena"
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
	fileResp, err := DecompressAndEncryptUserKubeConfigsWorkload(c)
	if err != nil {
		log.Err(err).Msg("CreateOrUpdateKubeConfig: DecompressUserKubeConfigsWorkload")
		return err
	}
	err = EncAndUpload(c.Request().Context(), ou.OrgID, fileResp, AgeEnc)
	if err != nil {
		log.Err(err).Msg("CreateOrUpdateKubeConfig: EncAndUpload")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, "ok")
}

func EncAndUpload(ctx context.Context, orgID int, in bytes.Buffer, ageEnc encryption.Age) error {
	fn := fmt.Sprintf("%d.kube.tar.gz", orgID)
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
	uploader := s3uploader.NewS3ClientUploader(athena.AthenaS3Manager)
	err = uploader.UploadFromInMemFs(ctx, p, input, inMemFsEnc)
	if err != nil {
		log.Err(err).Msg("CreateOrUpdateKubeConfig: Upload")
		return err
	}
	inMemFsK8sConfig := auth_startup.ExtClustersRunDigitalOceanS3BucketObjAuthProcedure(ctx, orgID, KeysCfg)
	k := zeus_core.K8Util{}
	k.ConnectToK8sFromInMemFsCfgPath(inMemFsK8sConfig)
	rawCfg, err := k.GetRawConfigs()
	if err != nil {
		log.Err(err).Msg("CreateOrUpdateKubeConfig: GetRawConfigs")
		return err
	}
	m := make(map[string]string)
	for ctxName, ai := range rawCfg.Clusters {
		if strings.Contains(ai.Server, "aws") {
			m[ctxName] = "aws"
			fmt.Println("aws command found")
			continue
		}
		if strings.Contains(ai.Server, "digtalocean") {
			m[ctxName] = "do"
			fmt.Println("digital ocean command found")
			continue
		}
		if strings.Contains(ai.Server, "ovh") {
			fmt.Println("ovh server found")
			m[ctxName] = "ovh"
			continue
		}
		if strings.Contains(ai.Server, "gke") || strings.Contains(ctxName, "gke") || strings.Contains(ctxName, "gcp") {
			m[ctxName] = "gcp"
			fmt.Println("gcp command found")
			continue
		}
	}
	var kcf []ext_clusters.ExtClusterConfig
	for ctxVal, cpName := range m {
		kcf = append(kcf, ext_clusters.ExtClusterConfig{
			CloudProvider: cpName,
			Context:       ctxVal,
			ContextAlias:  ctxVal,
			Env:           "none",
		})
	}
	err = ext_clusters.InsertOrUpdateExtClusterConfigs(ctx, org_users.NewOrgUserWithID(orgID, 0), kcf)
	if err != nil {
		log.Err(err).Msg("CreateOrUpdateKubeConfig: InsertOrUpdateExtClusterConfigs")
		return err
	}
	return err
}
