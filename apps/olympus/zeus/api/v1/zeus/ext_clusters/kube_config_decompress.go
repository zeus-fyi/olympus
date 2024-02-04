package zeus_v1_clusters_api

import (
	"bytes"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/transformations"
	"github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/memfs"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func DecompressUserKubeConfigsWorkload(c echo.Context) (chart_workload.TopologyBaseInfraWorkload, error) {
	nk := chart_workload.TopologyBaseInfraWorkload{}
	file, err := c.FormFile("kubeconfig")
	if err != nil {
		log.Err(err).Msg("DecompressUserKubeConfigsWorkload: FormFile")
		return nk, err
	}
	src, err := file.Open()
	if err != nil {
		log.Err(err).Msg("DecompressUserKubeConfigsWorkload: file.Open()")
		return nk, err
	}
	defer src.Close()
	in := bytes.Buffer{}
	if _, err = io.Copy(&in, src); err != nil {
		log.Err(err).Msg("DecompressUserKubeConfigsWorkload: Copy")
		return nk, err
	}
	nk, err = UnGzipKubeConfig(&in)
	if err != nil {
		log.Err(err).Interface("kubernetesWorkload", nk).Msg("DecompressUserKubeConfigsWorkload: UnGzipKubeConfig")
		return nk, err
	}
	return nk, err
}
func UnGzipKubeConfig(in *bytes.Buffer) (chart_workload.TopologyBaseInfraWorkload, error) {
	yr := transformations.YamlFileIO{}
	p := filepaths.Path{DirIn: "/tmp", DirOut: "/tmp", FnIn: "kubeconfig.tar.gz"}
	m := memfs.NewMemFs()
	err := m.MakeFileIn(&p, in.Bytes())
	if err != nil {
		return yr.TopologyBaseInfraWorkload, err
	}
	p.DirOut = "/kubeconfig"
	comp := compression.NewCompression()
	err = comp.UnGzipFromInMemFsOutToInMemFS(&p, m)
	if err != nil {
		return yr.TopologyBaseInfraWorkload, err
	}

	p.DirIn = "/kubeconfig"
	//err = yr.ReadK8sWorkloadInMemFsDir(p, m)
	return yr.TopologyBaseInfraWorkload, err
}
