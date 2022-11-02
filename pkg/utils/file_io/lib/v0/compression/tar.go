package compression

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func TarFolder(p structs.Path) error {
	files, err := ioutil.ReadDir(p.DirIn)
	if err != nil {
		return err
	}
	tarfile, err := os.Create(p.Fn)
	if err != nil {
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
		file, err := os.Open(p.DirIn + string(filepath.Separator) + fileInfo.Name())
		if err != nil {
			return err
		}
		defer file.Close()
		header := new(tar.Header)
		header.Name = file.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()
		err = tarfileW.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tarfileW, file)
		if err != nil {
			return err
		}
	}
	return err
}
