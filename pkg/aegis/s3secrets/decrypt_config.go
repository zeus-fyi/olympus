package s3secrets

import "github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"

func (s *S3Secrets) DecryptAndUnGzip(p *structs.Path) error {
	err := s.Age.DecryptToFile(p)
	if err != nil {
		return err
	}
	err = s.GzipDecompress(p)
	return err
}
