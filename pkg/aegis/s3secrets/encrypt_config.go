package s3secrets

import (
	"errors"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (s *S3Secrets) GzipAndEncrypt(p *structs.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	err := s.CreateTarGzipArchiveDir(p)
	if err != nil {
		return err
	}

	err = s.Age.Encrypt(p)
	return err
}
