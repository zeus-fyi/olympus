package s3secrets

import (
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (s *S3Secrets) DecryptAndUnGzipToInMemFs(p *structs.Path, unzipDir string) error {
	err := s.Age.DecryptToMemFsFile(p, s.MemFS)
	if err != nil {
		log.Err(err).Msgf("DecryptAndUnGzipToInMemFs, DecryptToMemFsFile %s", p.FileInPath())
		return err
	}

	// fn in is now the unencrypted version, so fn.out -> fn.in
	p.FnIn = p.FnOut
	p.DirOut = unzipDir
	err = s.UnGzipFromInMemFsOutToInMemFS(p, s.MemFS)
	if err != nil {
		log.Err(err).Msgf("DecryptAndUnGzipToInMemFs, UnGzipFromInMemFsOutToInMemFS %s", p.FileOutPath())
		return err
	}
	return err
}
