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
	return in, err
}
