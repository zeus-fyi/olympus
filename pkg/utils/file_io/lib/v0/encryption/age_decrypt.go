package encryption

import (
	"io"
	"os"

	"filippo.io/age"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func Decrypt(p structs.Path, privateKey string) error {
	identity, err := age.ParseX25519Identity(privateKey)
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
