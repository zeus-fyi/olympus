package zeus_v1_clusters_api

import (
	"bytes"
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	s3uploader "github.com/zeus-fyi/olympus/datastores/s3/upload"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
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
	err = EncAndUpload(c.Request().Context(), fileResp, encryption.Age{})
	if err != nil {
		log.Err(err).Msg("CreateOrUpdateKubeConfig: EncAndUpload")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, "ok")
}

func EncAndUpload(ctx context.Context, in bytes.Buffer, ageEnc encryption.Age) error {
	fn := "kube-temp.tar.gz.age"
	bucketName := "zeus-fyi"
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fn),
	}
	p := filepaths.Path{}
	inMemFsEnc := memfs.NewMemFs()
	err := ageEnc.EncryptItem(inMemFsEnc, &p, in.Bytes())
	if err != nil {
		log.Err(err).Msg("CreateOrUpdateKubeConfig: EncryptItem")
		return err
	}
	uploader := s3uploader.NewS3ClientUploader(athena.AthenaS3Manager)
	err = uploader.Upload(ctx, p, input)
	if err != nil {
		log.Err(err).Msg("CreateOrUpdateKubeConfig: Upload")
		return err
	}
	return err
}
