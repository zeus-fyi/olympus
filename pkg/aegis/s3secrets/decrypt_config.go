package s3secrets

import "github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"

func (s *S3Secrets) UnGzipAndDecrypt(p *structs.Path) error {
	err := s.UnGzip(p)
	if err != nil {
		return err
	}

	err = s.Age.Decrypt(p)
	return err
}
