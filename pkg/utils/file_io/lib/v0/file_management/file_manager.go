package file_management

import (
	"io/ioutil"
	"os"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/logging"
)

type FileManagerLib struct {
	Log logging.Logger
}

func (l *FileManagerLib) CreateFile(p structs.Path, data []byte) error {
	// make path if it doesn't exist
	if _, err := os.Stat(p.FileOutPath()); os.IsNotExist(err) {
		_ = os.MkdirAll(p.DirOut, 0700) // Create your dir
	}
	err := ioutil.WriteFile(p.FileOutPath(), data, 0644)
	return err
}

func (l *FileManagerLib) CreateV2FileOut(p structs.Path, data []byte) error {
	// make path if it doesn't exist
	if _, err := os.Stat(p.V2FileOutPath()); os.IsNotExist(err) {
		_ = os.MkdirAll(p.DirOut, 0700) // Create your dir
	}
	err := ioutil.WriteFile(p.V2FileOutPath(), data, 0644)
	return err
}

// OpenFile requires you to know that you need to close this
func (l *FileManagerLib) OpenFile(p structs.Path) (*os.File, error) {
	f, err := os.OpenFile(p.FileInPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	return f, l.Log.ErrHandler(err)
}

func (l *FileManagerLib) DeleteFile(p structs.Path) error {
	err := os.Remove(p.FileInPath())
	return l.Log.ErrHandler(err)
}
