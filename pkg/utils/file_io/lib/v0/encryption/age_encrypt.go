package encryption

import (
	"io"
	"os"

	"filippo.io/age"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/readers"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func Encrypt(p structs.Path, publicKey string) error {
	recipient, err := age.ParseX25519Recipient(publicKey)
	if err != nil {
		return err
	}

	outFile, err := os.Create(p.Fn + ".age")
	if err != nil {
		return err
	}
	defer outFile.Close()
	rl := readers.ReaderLib{}
	bytesToEncrypt := rl.ReadFile(p)

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
	return err
}
