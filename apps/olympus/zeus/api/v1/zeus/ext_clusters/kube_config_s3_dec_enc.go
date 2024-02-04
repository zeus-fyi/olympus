package zeus_v1_clusters_api

import (
	"bytes"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func DecompressAndEncryptUserKubeConfigsWorkload(c echo.Context) (bytes.Buffer, error) {
	in := bytes.Buffer{}
	file, err := c.FormFile("kubeconfig")
	if err != nil {
		log.Err(err).Msg("DecompressUserKubeConfigsWorkload: FormFile")
		return in, err
	}
	src, err := file.Open()
	if err != nil {
		log.Err(err).Msg("DecompressUserKubeConfigsWorkload: file.Open()")
		return in, err
	}
	defer src.Close()
	if _, err = io.Copy(&in, src); err != nil {
		log.Err(err).Msg("DecompressUserKubeConfigsWorkload: Copy")
		return in, err
	}
	//log.Info().Int("bytes", in.Len()).Msg("DecompressUserKubeConfigsWorkload: Copy")
	return in, err
}

//func UnGzipKubeConfig(in *bytes.Buffer) ([]byte, error) {
//	p := filepaths.Path{DirIn: "/tmp", DirOut: "/tmp", FnIn: "kubeconfig.tar.gz"}
//	m := memfs.NewMemFs()
//	err := m.MakeFileIn(&p, in.Bytes())
//	if err != nil {
//		return nil, err
//	}
//	p.DirOut = "/kubeconfig"
//	comp := compression.NewCompression()
//	err = comp.UnGzipFromInMemFsOutToInMemFS(&p, m)
//	if err != nil {
//		return nil, err
//	}
//
//	p.DirIn = "/kubeconfig"
//	return nil, err
//}
