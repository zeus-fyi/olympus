package encryption

import (
	"errors"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (a *Age) DecryptToMemFsFile(p *structs.Path, fs memfs.MemFS) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	out, err := a.decrypt(p)
	if err != nil {
		return err
	}

	err = fs.MakeFile(p, out.Bytes())
	return err
}
