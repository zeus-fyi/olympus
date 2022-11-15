package compression

import (
	"archive/tar"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (c *Compression) TarFolder(p *structs.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}

	files, err := ioutil.ReadDir(p.DirIn)
	if err != nil {
		return err
	}
	tarfile, err := os.Create(p.FnIn)
	if err != nil {
		log.Err(err).Msg("Compression: TarFolder, os.Create(p.FnIn)")
		return err
	}
	defer tarfile.Close()
	var fileW io.WriteCloser = tarfile
	//Tar file writter
	tarfileW := tar.NewWriter(fileW)
	defer tarfileW.Close()
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}
		file, ferr := os.Open(p.DirIn + string(filepath.Separator) + fileInfo.Name())
		if ferr != nil {
			log.Err(ferr).Msg("Compression: TarFolder, os.Open(p.DirIn + string(filepath.Separator) + fileInfo.Name())")
			return ferr
		}
		defer file.Close()
		header := new(tar.Header)
		header.Name = file.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()
		err = tarfileW.WriteHeader(header)
		if err != nil {
			log.Err(err).Msg("Compression: TarFolder, tarfileW.WriteHeader(header)")
			return err
		}
		_, err = io.Copy(tarfileW, file)
		if err != nil {
			log.Err(err).Msg("Compression: TarFolder, io.Copy(tarfileW, file)")
			return err
		}
	}
	return err
}
