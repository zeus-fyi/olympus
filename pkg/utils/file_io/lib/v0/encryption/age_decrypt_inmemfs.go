package encryption

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"filippo.io/age"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

func (a *Age) DecryptToMemFsFile(p *filepaths.Path, fs memfs.MemFS) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	out, err := a.decryptFromInMemFS(p, fs)
	if err != nil {
		log.Err(err).Msgf("DecryptToMemFsFile, decryptFromInMemFS %s", p.FileInPath())
		return err
	}

	err = fs.MakeFileOut(p, out.Bytes())
	if err != nil {
		log.Err(err).Msgf("DecryptToMemFsFile, MakeFileOut %s", p.FileOutPath())
		return err
	}
	return err
}

func (a *Age) decryptFromInMemFS(p *filepaths.Path, fs memfs.MemFS) (*bytes.Buffer, error) {
	out := &bytes.Buffer{}

	if p == nil {
		return out, errors.New("need to include a path")
	}
	identity, err := age.ParseX25519Identity(a.agePrivateKey)
	if err != nil {
		log.Err(err).Msg("Age, decryptFromInMemFS")
		return out, err
	}
	f, err := fs.Open(p.FileInPath())
	if err != nil {
		log.Err(err).Msgf("Age, decryptFromInMemFS, fs.Open(p.FileInPath()) %s", p.FileInPath())
		return out, err
	}
	defer f.Close()
	r, err := age.Decrypt(f, identity)
	if err != nil {
		log.Err(err).Msg("Age, decryptFromInMemFS, age.Decrypt")
		return out, err
	}
	p.FnOut, _, _ = strings.Cut(p.FnIn, ".age")
	if _, cerr := io.Copy(out, r); cerr != nil {
		log.Err(cerr).Msg("Age, decryptFromInMemFS, io.RsyncBucket")
		return out, cerr
	}
	return out, err
}
