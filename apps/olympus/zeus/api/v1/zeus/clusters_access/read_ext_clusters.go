package zeus_v1_clusters_api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_cluster_configs "github.com/zeus-fyi/olympus/pkg/hestia/cluster_configs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

type ReadPrivateClustersRequest struct {
}

func ReadExtKubeConfigsHandler(c echo.Context) error {
	request := new(ReadPrivateClustersRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ReadExtKubeConfig(c)
}

func (t *ReadPrivateClustersRequest) ReadExtKubeConfig(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	if ou.OrgID == 0 {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	extCfgs, err := authorized_clusters.SelectAuthedClusterConfigsByOrgID(c.Request().Context(), ou)
	if err != nil {
		log.Err(err).Msg("ReadExtKubeConfig: SelectExtClusterConfigsByOrgID")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	m := make(map[string]authorized_clusters.K8sClusterConfig)
	for _, cv := range extCfgs {
		m[fmt.Sprintf("%s-%s-%s", cv.CloudProvider, cv.Context, cv.Region)] = cv
	}
	cfgsRaw, err := hestia_cluster_configs.GetExtClusterConfigs(c.Request().Context(), ou)
	if err != nil {
		log.Err(err).Msg("ReadExtKubeConfig: GetExtClusterConfigs")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	for _, cfg := range cfgsRaw {
		if _, ok := m[fmt.Sprintf("%s-%s-%s", cfg.CloudProvider, cfg.Context, cfg.Region)]; !ok {
			inCmp, cerr := compression.GzipDirectoryToMemoryFS(cfg.Path, cfg.InMemFsKubeConfig)
			if cerr != nil {
				log.Err(cerr).Msg("ReadExtKubeConfig: GzipDirectoryToMemoryFS")
				return c.JSON(http.StatusInternalServerError, nil)
			}
			if inCmp == nil {
				log.Err(cerr).Msg("ReadExtKubeConfig: GzipDirectoryToMemoryFS")
				return c.JSON(http.StatusInternalServerError, nil)
			}
			if err != nil {
				log.Err(err).Msg("ReadExtKubeConfig: buf.Read")
				return c.JSON(http.StatusInternalServerError, nil)
			}
			err = EncAndUpload(c.Request().Context(), ou.OrgID, *inCmp, zeus.AgeEnc, cfg)
			if err != nil {
				log.Err(err).Msg("ReadExtKubeConfig: EncAndUpload")
				return c.JSON(http.StatusInternalServerError, nil)
			}
			extCfgs = append(extCfgs, cfg)
		}
	}
	return c.JSON(http.StatusOK, extCfgs)
}
