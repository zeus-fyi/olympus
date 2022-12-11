package zeus

import (
	"bytes"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
)

func DecompressUserInfraWorkload(c echo.Context) (chart_workload.TopologyBaseInfraWorkload, error) {
	nk := chart_workload.TopologyBaseInfraWorkload{}
	file, err := c.FormFile("chart")
	if err != nil {
		log.Err(err).Msg("DecompressUserInfraWorkload: FormFile")
		return nk, err
	}
	src, err := file.Open()
	if err != nil {
		log.Err(err).Msg("DecompressUserInfraWorkload: file.Open()")
		return nk, err
	}
	defer src.Close()
	in := bytes.Buffer{}
	if _, err = io.Copy(&in, src); err != nil {
		log.Err(err).Msg("DecompressUserInfraWorkload: RsyncBucket")
		return nk, err
	}
	nk, err = UnGzipK8sChart(&in)
	if err != nil {
		log.Err(err).Interface("kubernetesWorkload", nk).Msg("DecompressUserInfraWorkload: UnGzipK8sChart")
		return nk, err
	}
	return nk, err
}
