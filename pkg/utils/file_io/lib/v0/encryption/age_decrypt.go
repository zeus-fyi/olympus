package encryption

import (
	"errors"
	"io"
	"os"

	"filippo.io/age"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (a *Age) Decrypt(p *structs.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	identity, err := age.ParseX25519Identity(a.agePrivateKey)
	if err != nil {
		return err
	}
	f, err := os.Open(p.Fn)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := age.Decrypt(f, identity)
	if err != nil {
		return err
	}
	outFile, err := os.Create(p.FnOut)
	if err != nil {
		return err
	}
	if _, cerr := io.Copy(outFile, r); cerr != nil {
		return cerr
	}
	return err
}
