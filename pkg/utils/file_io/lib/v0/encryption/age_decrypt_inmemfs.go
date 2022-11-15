package encryption

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"filippo.io/age"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (a *Age) DecryptToMemFsFile(p *structs.Path, fs memfs.MemFS) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	out, err := a.decryptFromInMemFS(p, fs)
	if err != nil {
		return err
	}

	err = fs.MakeFile(p, out.Bytes())
	return err
}

func (a *Age) decryptFromInMemFS(p *structs.Path, fs memfs.MemFS) (*bytes.Buffer, error) {
	out := &bytes.Buffer{}

	if p == nil {
		return out, errors.New("need to include a path")
	}
	identity, err := age.ParseX25519Identity(a.agePrivateKey)
	if err != nil {
		return out, err
	}
	f, err := fs.Open(p.FnIn)
	if err != nil {
		return out, err
	}
	defer f.Close()
	r, err := age.Decrypt(f, identity)
	if err != nil {
		return out, err
	}
	p.FnOut, _, _ = strings.Cut(p.FnIn, ".age")
	if _, cerr := io.Copy(out, r); cerr != nil {
		return out, cerr
	}
	p.FnIn = p.FnOut
	return out, err
}
