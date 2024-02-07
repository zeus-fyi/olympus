package zeus_v1_clusters_api

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type ReadAuthorizedClustersRequest struct {
}

func ReadAuthorizedClustersRequestHandler(c echo.Context) error {
	request := new(ReadAuthorizedClustersRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Read(c)
}

func (t *ReadAuthorizedClustersRequest) Read(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	if ou.OrgID == 0 {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	extCfgs, err := authorized_clusters.SelectAuthedClusterConfigsByOrgID(c.Request().Context(), ou)
	if err != nil {
		log.Err(err).Msg("ReadExtKubeConfig: SelectAuthedClusterConfigsByOrgID")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, extCfgs)
}

func CheckKubeConfig(ctx context.Context, ou org_users.OrgUser, p authorized_clusters.K8sClusterConfig) error {
	inKey := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String(p.GetHashedKey(ou.OrgID)),
	}
	inMemFsK8sConfig := auth_startup.ExtClustersRunDigitalOceanS3BucketObjAuthProcedure(ctx, KeysCfg, inKey)
	k := zeus_core.K8Util{}
	err := k.ConnectToK8sFromInMemFsCfgPathOrErr(inMemFsK8sConfig)
	if err != nil {
		log.Err(err).Msg("CheckKubeConfig: ConnectToK8sFromInMemFsCfgPathOrErr")
	}
	return err
}
