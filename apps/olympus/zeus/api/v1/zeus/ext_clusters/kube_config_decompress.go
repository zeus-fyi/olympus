package zeus_v1_clusters_api

import (
	"bytes"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/memfs"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func DecompressUserKubeConfigsWorkload(c echo.Context) ([]byte, error) {
	file, err := c.FormFile("kubeconfig")
	if err != nil {
		log.Err(err).Msg("DecompressUserKubeConfigsWorkload: FormFile")
		return nil, err
	}
	src, err := file.Open()
	if err != nil {
		log.Err(err).Msg("DecompressUserKubeConfigsWorkload: file.Open()")
		return nil, err
	}
	defer src.Close()
	in := bytes.Buffer{}
	if _, err = io.Copy(&in, src); err != nil {
		log.Err(err).Msg("DecompressUserKubeConfigsWorkload: Copy")
		return nil, err
	}
	log.Info().Int("bytes", in.Len()).Msg("DecompressUserKubeConfigsWorkload: Copy")

	b, err := UnGzipKubeConfig(&in)
	if err != nil {
		log.Err(err).Msg("DecompressUserKubeConfigsWorkload: UnGzipKubeConfig")
		return nil, err
	}
	return b, err
}
func UnGzipKubeConfig(in *bytes.Buffer) ([]byte, error) {
	p := filepaths.Path{DirIn: "/tmp", DirOut: "/tmp", FnIn: "kubeconfig.tar.gz"}
	m := memfs.NewMemFs()
	err := m.MakeFileIn(&p, in.Bytes())
	if err != nil {
		return nil, err
	}
	p.DirOut = "/kubeconfig"
	comp := compression.NewCompression()
	err = comp.UnGzipFromInMemFsOutToInMemFS(&p, m)
	if err != nil {
		return nil, err
	}

	p.DirIn = "/kubeconfig"
	return nil, err
}
