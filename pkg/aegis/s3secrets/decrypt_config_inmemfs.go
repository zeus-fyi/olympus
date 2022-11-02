package s3secrets

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (s *S3Secrets) DecryptAndUnGzipToInMemFs(p *structs.Path, unzipDir string) error {
	err := s.Age.DecryptToMemFsFile(p, s.MemFS)
	if err != nil {
		return err
	}
	p.DirOut = unzipDir
	err = s.UnGzipFromInMemFsOutToInMemFS(p, s.MemFS)
	return err
}
