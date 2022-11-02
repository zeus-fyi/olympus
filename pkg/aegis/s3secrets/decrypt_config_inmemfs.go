package s3secrets

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (s *S3Secrets) DecryptAndUnGzipToInMemFs(p *structs.Path, unzipDir string, fs memfs.MemFS) error {
	err := s.Age.DecryptToMemFsFile(p, fs)
	if err != nil {
		return err
	}
	p.DirOut = unzipDir
	err = s.UnGzipFromInMemFsOutToInMemFS(p, fs)
	return err
}
