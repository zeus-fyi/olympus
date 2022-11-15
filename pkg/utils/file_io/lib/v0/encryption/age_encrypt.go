package encryption

import (
	"errors"
	"io"
	"os"

	"filippo.io/age"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/readers"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (a *Age) Encrypt(p *structs.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	recipient, err := age.ParseX25519Recipient(a.agePublicKey)
	if err != nil {
		return err
	}
	p.FnOut = p.FnIn + ".age"
	outFile, err := os.Create(p.FnOut)
	if err != nil {
		return err
	}
	defer outFile.Close()
	rl := readers.ReaderLib{}
	bytesToEncrypt := rl.ReadFilePathPtr(p)

	w, err := age.Encrypt(outFile, recipient)
	if err != nil {
		return err
	}
	if _, werr := w.Write(bytesToEncrypt); werr != nil {
		return werr
	}
	_, err = io.Copy(w, outFile)
	if err != nil {
		return err
	}
	if cerr := w.Close(); cerr != nil {
		return cerr
	}
	p.FnIn = p.FnOut
	return err
}
