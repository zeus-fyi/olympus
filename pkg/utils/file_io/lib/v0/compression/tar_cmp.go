package compression

import (
	"archive/tar"
	"errors"
	"io/fs"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func (c *Compression) TarCompress(p *filepaths.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	p.FnOut = p.FnIn + ".tar"
	out, err := os.Create(p.FileOutPath())
	if err != nil {
		log.Err(err).Msg("Compression: TarCompress, os.Create(p.FileOutPath()")
		return err
	}
	defer out.Close()

	tw := tar.NewWriter(out)
	defer tw.Close()
	fileSystem := os.DirFS(p.DirIn)
	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			aerr := addToArchive(p, tw, path)
			if aerr != nil {
				return aerr
			}
		}
		return nil
	})

	p.DirIn = p.DirOut
	p.FnIn = p.FnOut
	return err
}
